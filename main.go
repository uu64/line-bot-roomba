package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var logger *botLogger

type botLogger struct {
	stdout io.Writer
	stderr io.Writer
}

func (l *botLogger) info(message string) {
	fmt.Fprintln(l.stdout, message)
}

func (l *botLogger) error(message string) {
	fmt.Fprintln(l.stderr, message)
}

func newLogger() *botLogger {
	return &botLogger{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

type iftttReqBody struct {
	Event string `json:"event"`
}

func init() {
	// logging setting
	logger = newLogger()

	// line bot settings
	var err error
	token := os.Getenv("LINE_BOT_TOKEN")
	secret := os.Getenv("LINE_BOT_SECRET")
	bot, err = linebot.New(secret, token)
	if err != nil {
		logger.error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

func push(text string) error {
	to := os.Getenv("LINE_BOT_PRIVATE_ID")
	message := linebot.NewTextMessage(text)
	if _, err := bot.PushMessage(to, message).Do(); err != nil {
		return err
	}
	return nil
}

func reply(text string, replyToken string) error {
	message := linebot.NewTextMessage(text)
	if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
		return err
	}
	return nil
}

func iftttHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bot/ifttt" {
		http.NotFound(w, r)
		logger.error("invalid path")
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		logger.error("invalid method")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		logger.error("invalid header")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.error(fmt.Sprintf("failed to read body: %+v", err))
		return
	}

	var body iftttReqBody
	if err = json.Unmarshal(b, &body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.error(fmt.Sprintf("failed to unmarshal: %+v", err))
		return
	}

	if body.Event == "finish-cleaning" {
		if err := push("掃除おわった"); err != nil {
			logger.error(fmt.Sprintf("failed to send a push messsage: %+v", err))
		}
	}

	if body.Event == "be-stuck" {
		if err := push("たす...け...て......"); err != nil {
			logger.error(fmt.Sprintf("failed to send a push messsage: %+v", err))
		}
	}
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bot/webhook" {
		http.NotFound(w, r)
		logger.error("invalid path")
		return
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			logger.error(fmt.Sprintf("invalid signature: %+v", err))
		} else {
			w.WriteHeader(500)
			logger.error(fmt.Sprintf("failed to parse webhook request: %+v", err))
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				nicknames := []string{
					"ルンちゃん",
					"るんちゃん",
					"ルンさん",
					"るんさん",
				}
				for _, nickname := range nicknames {
					if strings.Contains(message.Text, nickname) {
						if err := reply("ほい", event.ReplyToken); err != nil {
							logger.error(fmt.Sprintf("failed to send a reply messsage: %+v", err))
						}
						return
					}
				}
			}
		}
	}
}

func main() {
	// initialize handler
	http.HandleFunc("/bot/webhook", webhookHandler)
	http.HandleFunc("/bot/ifttt", iftttHandler)

	// initialize server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.info(fmt.Sprintf("Defaulting to port %s", port))
	}

	logger.info(fmt.Sprintf("Listening on port %s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
