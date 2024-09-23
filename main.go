package main

import (
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
		// TODO: implement game progress
		// tickGameProgress(c)
		time.Sleep(1 * time.Minute)
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
			log.Printf("error during reply check: %v\n", err)
			continue
		} else if exists {
			// TODO: check if move changed
			log.Printf("reply already exists: id=%d\n", n.Item.Id)
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

		var cId int
		parentId := n.Item.Id
		if cId, err = c.CreateComment(parentId, imgUrl); err != nil {
			log.Printf("error creating reply: %v\n", err)
			continue
		}

		n.Item.ParentId = 0
		if err = db.InsertItem(&n.Item); err != nil {
			log.Printf("error inserting item into db: %v: id=%d\n", err, n.Item.Id)
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
