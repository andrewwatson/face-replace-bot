// https://www.opsdash.com/blog/slack-bot-in-golang.html

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

var (
	token        string
	chrisifyPath string
	haarPath     string

	defaultRegion string
	accessKeyID   string
	secretKeyID   string
	s3Bucket      string

	sess *session.Session
	svc  *s3.S3
)

// var base_path = "/var/www/chrisbot.zikes.me/"
// var base_url = "http://chrisbot.zikes.me/"

func init() {
	token = os.Getenv("SLACK_TOKEN")
	if token == "" {
		panic(errors.New("SLACK_TOKEN must be provided"))
	}

	chrisifyPath = os.Getenv("CHRISIFY_PATH")
	haarPath = os.Getenv("HAAR_FILE")

	accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKeyID = os.Getenv("AWS_SECRET_ACCESS_KEY")

	s3Bucket = os.Getenv("S3_BUCKET_NAME")

	var err error
	sess, err = session.NewSession(&aws.Config{
		Region:      aws.String(defaultRegion),
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		panic(err.Error())
	}

	svc = s3.New(sess)

}

func main() {

	ws, id := slackConnect(token)
	fmt.Println("bot ready to replace faces")

	for {

		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && m.SubType == "file_share" && strings.Contains(m.Text, "<@"+id+">") {
			go func(m Message) {
				var channel string
				json.Unmarshal(m.Channel, &channel)
				file := SaveTempFile(GetFile(m.File))
				chrisd, err := Chrisify(file)
				if err != nil {
					log.Printf("Error during Face Replace: %s", err.Error())
				}

				unique := uuid.New()

				object, err := svc.PutObject(&s3.PutObjectInput{
					Bucket:      aws.String(s3Bucket),
					Key:         aws.String(unique.String()),
					Body:        bytes.NewReader(chrisd),
					ContentType: aws.String("image/png"),
					ACL:         aws.String("public-read"),
				})

				if err != nil {
					log.Printf("Error During Upload: %s", err.Error())
				} else {
					log.Printf("Object Created: %s", object)
				}

				postMessage(ws, map[string]string{
					"type":    "message",
					"text":    "https://s3.amazonaws.com/makeandbuild-kenbot-results/" + unique.String(),
					"channel": channel,
				})

				// defer os.Remove(file)
			}(m)
		}
	}

	// log.Fatal(errors.New("We Exited the For Loop!"))

}
