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
	c = sn.GetClient()
	// TODO: fetch our id from SN API
	// prod: 25176 | local: 21858
	meId = 25176
)

func main() {

	for {
		tickGameStart(c)
		tickGameProgress(c)
		time.Sleep(15 * time.Second)
	}
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

	log.Printf("fetched %d mention(s)\n", len(mentions))

	for _, n := range mentions {

		// we only care about current notifications
		x := time.Now().Add(-30 * time.Second)
		if n.Item.CreatedAt.Before(x) {
			log.Printf("ignoring old mention %d\n", n.Item.Id)
			continue
		}

		if exists, err := db.ItemHasReply(n.Item.Id, meId); err != nil {
			log.Printf("failed to check for existing reply to game start in item %d: %v\n", n.Item.Id, err)
			continue
		} else if exists {
			// TODO: check if move changed
			log.Printf("reply to game start in item %d already exists\n", n.Item.Id)
			continue
		}

		if err = handleGameStart(&n.Item); err != nil {

			// don't reply to mentions that we failed to parse as a game start
			// to support unrelated mentions
			if err.Error() == "failed to parse game start" {
				log.Printf("ignoring error for item %d: %v\n", n.Item.Id, err)
				return
			}

			if _, err2 := createComment(n.Item.Id, fmt.Sprintf("`%v`", err)); err2 != nil {
				log.Printf("failed to reply with error to item %d: %v\n", n.Item.Id, err2)
			} else {
				log.Printf("replied to game start in item %d with error: %v\n", n.Item.Id, err)
			}
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

		// we only care about current notifications
		x := time.Now().Add(-30 * time.Second)
		if n.Item.CreatedAt.Before(x) {
			log.Printf("ignoring old reply %d\n", n.Item.Id)
			continue
		}

		if exists, err := db.ItemHasReply(n.Item.Id, meId); err != nil {
			log.Printf("failed to check for existing reply to game update in item %d: %v\n", n.Item.Id, err)
			continue
		} else if exists {
			// TODO: check if move changed
			log.Printf("reply to game update in item %d already exists\n", n.Item.Id)
			continue
		}

		if err = handleGameProgress(&n.Item); err != nil {

			if err.Error() == "failed to parse game update" {
				log.Printf("ignoring error for item %d: %v\n", n.Item.Id, err)
				return
			}

			if _, err2 := createComment(n.Item.Id, fmt.Sprintf("`%v`", err)); err2 != nil {
				log.Printf("failed to reply with error to item %d: %v\n", n.Item.Id, err2)
			} else {
				log.Printf("replied to game start in item %d with error: %v\n", n.Item.Id, err)
			}
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

	// reply with algebraic notation and image
	res = strings.Trim(fmt.Sprintf("%s\n\n%s", b.AlgebraicNotation(), imgUrl), " ")
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
		if item.User.Id == meId {
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

		if !strings.HasPrefix(line, "@chess") {
			continue
		}

		return strings.Trim(strings.ReplaceAll(line, "@chess", ""), " "), nil
	}

	return "", errors.New("failed to parse game start")
}

func parseGameProgress(input string) (string, error) {
	lines := strings.Split(input, "\n")
	words := strings.Split(input, " ")

	if len(lines) == 1 && len(words) == 1 {
		return strings.Trim(input, " "), nil
	}

	for _, line := range strings.Split(input, "\n") {
		if !strings.Contains(line, "@chess") {
			continue
		}

		return strings.Trim(strings.ReplaceAll(line, "@chess", ""), " "), nil
	}

	return "", errors.New("failed to parse game update")
}
