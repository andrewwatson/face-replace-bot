// https://www.opsdash.com/blog/slack-bot-in-golang.html

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var token string
var chrisifyPath string
var haarPath string

// var base_path = "/var/www/chrisbot.zikes.me/"
// var base_url = "http://chrisbot.zikes.me/"

func init() {
	token = os.Getenv("SLACK_TOKEN")
	if token == "" {
		panic(errors.New("SLACK_TOKEN must be provided"))
	}

	chrisifyPath = os.Getenv("CHRISIFY_PATH")

	haarPath = os.Getenv("HAAR_FILE")
}

func main() {

	// if len(os.Args) != 4 {
	// 	fmt.Fprintf(os.Stderr, "usage: slackbot slack-bot-token /path/to/chrisify /path/to/haar\n")
	// 	os.Exit(1)
	// }
	//
	// token = os.Args[1]
	// chrisify = os.Args[2]
	// haar = os.Args[3]

	// start a websocket-based Real Time API session
	ws, id := slackConnect(token)
	fmt.Println("slackbot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && m.SubType == "file_share" && strings.Contains(m.Text, "<@"+id+">") {
			go func(m Message) {
				var channel string
				json.Unmarshal(m.Channel, &channel)
				// file := SaveTempFile(GetFile(m.File))
				// chrisd := Chrisify(file)
				// log.Printf("Uploading to %s", channel)
				// Upload(chrisd, channel)
				// url := SaveFile(chrisd)
				postMessage(ws, map[string]string{
					"type":    "message",
					"text":    "https://avatars.slack-edge.com/2017-02-24/145511248880_386a6bad513462a96741_48.png",
					"channel": channel,
				})

				// defer os.Remove(file)
			}(m)
		}
	}
}
