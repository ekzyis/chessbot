package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/ekzyis/chessbot/chess"
	"github.com/ekzyis/chessbot/db"
	"github.com/ekzyis/chessbot/sn"
)

var (
	c  = sn.GetClient()
	me *sn.User
)

func main() {

	for {
		updateMe()
		tickGameStart(c)
		tickGameProgress(c)
		time.Sleep(15 * time.Second)
	}
}

func updateMe() {
	var (
		oldMe         *sn.User
		err           error
		warnThreshold = 100
	)

	maybeWarn := func() {
		if me.Privates.Sats < warnThreshold {
			log.Printf("~~~ warning: low balance ~~~\n")
		}
	}

	if me == nil {
		// make sure first update is successful
		if me, err = c.Me(); err != nil {
			log.Fatalf("failed to fetch me: %v\n", err)
		}
		log.Printf("fetched me: id=%d name=%s balance=%d\n", me.Id, me.Name, me.Privates.Sats)
		maybeWarn()
		return
	}

	oldMe = me
	if me, err = c.Me(); err != nil {
		log.Printf("failed to update me: %v\n", err)
		me = oldMe
	} else {
		log.Printf("updated me. balance: %d\n", me.Privates.Sats)
	}

	maybeWarn()
}

func tickGameStart(c *sn.Client) {
	var (
		mentions []sn.Notification
		err      error
	)

	if mentions, err = c.Mentions(); err != nil {
		log.Printf("failed to fetch mentions: %v\n", err)
		return
	}

	log.Printf("fetched %d mentions\n", len(mentions))

	for _, n := range mentions {

		if !isRecent(n.Item.CreatedAt) {
			log.Printf("ignoring old mention %d\n", n.Item.Id)
			continue
		}

		if handled, err := alreadyHandled(n.Item.Id); err != nil {
			log.Printf("failed to check for existing reply to game start in item %d: %v\n", n.Item.Id, err)
			continue
		} else if handled {
			// TODO: check if move changed
			log.Printf("reply to game start in item %d already exists\n", n.Item.Id)
			continue
		}

		if err = handleGameStart(&n.Item); err != nil {
			handleError(&n.Item, err)
		} else {
			log.Printf("started new game via item %d\n", n.Item.Id)
		}
	}
}

func tickGameProgress(c *sn.Client) {
	var (
		replies []sn.Notification
		err     error
	)

	if replies, err = c.Replies(); err != nil {
		log.Printf("failed to fetch replies: %v\n", err)
		return
	}

	log.Printf("fetched %d replies\n", len(replies))

	for _, n := range replies {

		if !isRecent(n.Item.CreatedAt) {
			log.Printf("ignoring old reply %d\n", n.Item.Id)
			continue
		}

		if handled, err := alreadyHandled(n.Item.Id); err != nil {
			log.Printf("failed to check for existing reply to game update in item %d: %v\n", n.Item.Id, err)
			continue
		} else if handled {
			// TODO: check if move changed
			log.Printf("reply to game update in item %d already exists\n", n.Item.Id)
			continue
		}

		if parent, err := c.Item(n.Item.ParentId); err != nil {
			log.Printf("failed to fetch parent %d of %d\n", n.Item.ParentId, n.Item.Id)
			continue
		} else if parent.User.Id != me.Id {
			log.Printf("ignoring nested reply %d\n", n.Item.Id)
			continue
		}

		if err = handleGameProgress(&n.Item); err != nil {
			handleError(&n.Item, err)
		} else {
			log.Printf("updated game via item %d\n", n.Item.Id)
		}
	}
}

func handleGameStart(req *sn.Item) error {
	var (
		move   string
		b      *chess.Board
		imgUrl string
		res    string
		err    error
	)

	// Immediately save game start request to db so we can store our reply to it in case of error.
	// We set parentId to 0 such that parent_id will be NULL in the db and not hit foreign key constraints.
	req.ParentId = 0

	if err = db.InsertItem(req); err != nil {
		return fmt.Errorf("failed to insert item %d into db: %v\n", req.Id, err)
	}

	if move, err = parseGameStart(req.Text); err != nil {
		return err
	}

	// create board with initial move(s)
	if b, err = chess.NewGame(move); err != nil {
		if rand.Float32() > 0.99 {
			// easter egg error message
			return errors.New("Nice try, fed.")
		}
		return fmt.Errorf("failed to create new game from item %d: %v\n", req.Id, err)
	}

	// upload image of board
	if imgUrl, err = c.UploadImage(b.Image()); err != nil {
		return fmt.Errorf("failed to upload image for item %d: %v\n", req.Id, err)
	}

	// reply with algebraic notation, image and info
	infoMove := "e4"
	if len(b.Moves) > 0 {
		infoMove = "e5"
	}
	info := fmt.Sprintf("_A new chess game has been started!_\n\n"+
		"_Reply with a move like `%s` to continue the game. "+
		"See [here](https://stacker.news/chess#how-to-continue) for details._", infoMove)
	res = strings.Trim(fmt.Sprintf("%s\n\n%s\n\n%s", b.AlgebraicNotation(), imgUrl, info), " ")
	if _, err = createComment(req.Id, res); err != nil {
		return fmt.Errorf("failed to reply to item %d: %v\n", req.Id, err)
	}

	return nil
}

func handleGameProgress(req *sn.Item) error {
	var (
		thread []sn.Item
		b             = chess.NewBoard()
		move   string = strings.Trim(req.Text, " ")
		imgUrl string
		res    string
		err    error
	)

	// immediately save game update request to db so we can store our reply to it in case of error
	if err = db.InsertItem(req); err != nil {
		return fmt.Errorf("failed to insert item %d into db: %v\n", req.Id, err)
	}

	// fetch thread to reconstruct all moves so far
	if thread, err = db.GetThread(req.ParentId); err != nil {
		return fmt.Errorf("failed to fetch thread for item %d: %v\n", req.ParentId, err)
	}

	for _, item := range thread {
		if item.User.Id == me.Id {
			continue
		}

		var moves string
		if moves, err = parseGameProgress(item.Text); err != nil {
			return err
		}

		// parse and execute existing moves
		if err = b.Parse(moves); err != nil {
			return err
		}
	}

	// parse and execute new move

	if move, err = parseGameProgress(move); err != nil {
		return err
	}

	if err = b.Parse(move); err != nil {
		if rand.Float32() > 0.99 {
			// easter egg error message
			return errors.New("Nice try, fed.")
		}
		return err
	}

	// upload image of updated board
	if imgUrl, err = c.UploadImage(b.Image()); err != nil {
		return fmt.Errorf("failed to upload image for item %d: %v\n", req.Id, err)
	}

	// reply with algebraic notation and image
	res = strings.Trim(fmt.Sprintf("%s\n\n%s", b.AlgebraicNotation(), imgUrl), " ")
	if _, err = createComment(req.Id, res); err != nil {
		return fmt.Errorf("failed to reply to item %d: %v\n", req.Id, err)
	}

	return nil
}

func handleError(req *sn.Item, err error) {

	// don't reply to mentions that we failed to parse as a game start
	// to support unrelated mentions
	if err.Error() == "failed to parse game start" {
		log.Printf("ignoring error for item %d: %v\n", req.Id, err)
		return
	}

	if err.Error() == "failed to parse game update" {
		log.Printf("ignoring error for item %d: %v\n", req.Id, err)
		return
	}

	if _, err2 := createComment(req.Id, fmt.Sprintf("`%v`", err)); err2 != nil {
		log.Printf("failed to reply with error to item %d: %v\n", req.Id, err2)
	} else {
		log.Printf("replied to game start in item %d with error: %v\n", req.Id, err)
	}
}

func createComment(parentId int, text string) (*sn.Item, error) {
	var (
		commentId int
		err       error
	)

	if commentId, err = c.CreateComment(parentId, text); err != nil {
		return nil, fmt.Errorf("failed to reply to item %d: %v\n", parentId, err)
	}

	var comment *sn.Item
	if comment, err = c.Item(commentId); err != nil {
		return nil, fmt.Errorf("failed to fetch item %d: %v\n", commentId, err)
	}

	if err = db.InsertItem(comment); err != nil {
		return nil, fmt.Errorf("failed to insert item %d into db: %v\n", comment.Id, err)
	}

	return comment, nil
}

func parseGameStart(input string) (string, error) {
	for _, line := range strings.Split(input, "\n") {
		line = strings.Trim(line, " ")

		var found bool
		if line, found = strings.CutPrefix(line, "@chess"); !found {
			continue
		}

		return strings.Trim(line, " "), nil
	}

	return "", errors.New("failed to parse game start")
}

func parseGameProgress(input string) (string, error) {
	input = strings.Trim(input, " ")

	lines := strings.Split(input, "\n")
	words := strings.Split(input, " ")

	if len(lines) == 1 && len(words) == 1 {
		return strings.Trim(strings.ReplaceAll(input, "@chess", ""), " "), nil
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.Trim(line, " ")

		var found bool
		if line, found = strings.CutPrefix(line, "@chess"); !found {
			continue
		}

		return strings.Trim(line, " "), nil
	}

	return "", errors.New("failed to parse game update")
}

func isRecent(t time.Time) bool {
	x := time.Now().Add(-30 * time.Second)
	return t.After(x)
}

func alreadyHandled(id int) (bool, error) {
	return db.ItemHasReply(id, me.Id)
}
