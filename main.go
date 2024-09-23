package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ekzyis/chessbot/chess"
	"github.com/ekzyis/chessbot/db"
	"github.com/ekzyis/chessbot/sn"
)

var (
	c = sn.GetClient()
	// TODO: fetch our id from SN API
	meId = 21858
)

func main() {

	for {
		tickGameStart(c)
		tickGameProgress(c)
		time.Sleep(30 * time.Second)
	}
}

func tickGameStart(c *sn.Client) {
	var (
		mentions []sn.Notification
		err      error
	)

	if mentions, err = c.Mentions(); err != nil {
		log.Printf("mentions error: %v\n", err)
		return
	}

	log.Printf("fetched %d mention(s)\n", len(mentions))

	for _, n := range mentions {

		if exists, err := db.ItemHasReply(n.Item.Id, meId); err != nil {
			log.Printf("error during game start check: %v\n", err)
			continue
		} else if exists {
			// TODO: check if move changed
			log.Printf("game start already exists: id=%d\n", n.Item.Id)
			continue
		}

		move := strings.Trim(strings.ReplaceAll(n.Item.Text, "@chess", ""), " ")

		var b *chess.Board
		if b, err = chess.NewGame(move); err != nil {
			log.Printf("error creating new game: %v: id=%d\n", err, n.Item.Id)
			continue
		}

		img := b.Image()
		var imgUrl string
		if imgUrl, err = c.UploadImage(img); err != nil {
			log.Printf("error uploading image: %v\n", err)
			continue
		}

		n.Item.ParentId = 0
		if err = db.InsertItem(&n.Item); err != nil {
			log.Printf("error inserting item into db: %v: id=%d\n", err, n.Item.Id)
			continue
		}

		text := fmt.Sprintf("`%s`\n\n%s", b.AlgebraicNotation(), imgUrl)
		var cId int
		if cId, err = c.CreateComment(n.Item.Id, text); err != nil {
			log.Printf("error creating reply: %v\n", err)
			continue
		}

		var item *sn.Item
		if item, err = c.Item(cId); err != nil {
			log.Printf("error fetching item: %v: id=%d\n", cId)
			continue
		}
		if err = db.InsertItem(item); err != nil {
			log.Printf("error inserting item into db: %v: id=%d\n", err, item.Id)
			continue
		}

		log.Printf("started new game: id=%d\n", n.Item.Id)
	}
}

func tickGameProgress(c *sn.Client) {
	var (
		replies []sn.Notification
		err     error
	)

	if replies, err = c.Replies(); err != nil {
		log.Printf("replies error: %v\n", err)
		return
	}

	log.Printf("fetched %d replies\n", len(replies))

	for _, n := range replies {

		if exists, err := db.ItemHasReply(n.Item.Id, meId); err != nil {
			log.Printf("error during game update check: %v\n", err)
			continue
		} else if exists {
			// TODO: check if move changed
			log.Printf("game update already exists: id=%d\n", n.Item.Id)
			continue
		}

		var thread []sn.Item
		if thread, err = db.GetThread(n.Item.ParentId); err != nil {
			log.Printf("error fetching thread: %v: id=%d\n", err, n.Item.ParentId)
			continue
		}

		b := chess.NewBoard()
		// TODO: better parsing of moves in replies using regexp for example or enforce a specific format
		// since players might include more than just the move in their replies
		move := strings.Trim(n.Item.Text, " ")
		for _, item := range thread {
			if item.User.Id == meId {
				continue
			}
			move := strings.Trim(strings.ReplaceAll(item.Text, "@chess", ""), " ")
			if err = b.Move(move); err != nil {
				log.Printf("error moving piece: %v: id=%d\n", err, item.Id)
				break
			}
		}

		// continue with next reply if the thread loop failed
		if err != nil {
			continue
		}

		if err = b.Move(move); err != nil {
			log.Printf("error moving piece: %v: id=%d\n", err, n.Item.ParentId)
			continue
		}

		img := b.Image()
		var imgUrl string
		if imgUrl, err = c.UploadImage(img); err != nil {
			log.Printf("error uploading image: %v\n", err)
			continue
		}

		if err = db.InsertItem(&n.Item); err != nil {
			log.Printf("error inserting item into db: %v: id=%d\n", err, n.Item.Id)
			continue
		}

		text := fmt.Sprintf("`%s`\n\n%s", b.AlgebraicNotation(), imgUrl)
		var cId int
		if cId, err = c.CreateComment(n.Item.Id, text); err != nil {
			log.Printf("error creating reply: %v\n", err)
			continue
		}

		var item *sn.Item
		if item, err = c.Item(cId); err != nil {
			log.Printf("error fetching item: %v: id=%d\n", cId)
			continue
		}
		if err = db.InsertItem(item); err != nil {
			log.Printf("error inserting item into db: %v: id=%d\n", err, item.Id)
			continue
		}

		log.Printf("continued game: move=%s id=%d\n", move, cId)
	}
}
