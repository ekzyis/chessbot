package sn

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	snappy "github.com/ekzyis/snappy"
)

type Client = snappy.Client
type Notification = snappy.Notification
type Item = snappy.Item
type User = snappy.User

var (
	c *Client
)

func GetClient() *Client {
	loadEnv()

	if c == nil {
		c = snappy.NewClient(
			snappy.WithBaseUrl(os.Getenv("SN_BASE_URL")),
			snappy.WithApiKey(os.Getenv("SN_API_KEY")),
			snappy.WithMediaUrl(os.Getenv("SN_MEDIA_URL")),
		)
	}
	return c
}

func loadEnv() {
	var (
		f   *os.File
		s   *bufio.Scanner
		err error
	)

	if f, err = os.Open(".env"); err != nil {
		log.Fatalf("error opening .env: %v", err)
	}
	defer f.Close()

	s = bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, "=", 2)

		// Check if we have exactly 2 parts (key and value)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		} else {
			log.Fatalf(".env: invalid line: %s\n", line)
		}
	}

	// Check for errors during scanning
	if err = s.Err(); err != nil {
		fmt.Println("error scanning .env:", err)
	}
}
