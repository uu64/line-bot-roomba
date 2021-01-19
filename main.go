package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func init() {
	token := os.Getenv("LINE_BOT_TOKEN")
	secret := os.Getenv("LINE_BOT_SECRET")

	var err error
	bot, err = linebot.New(secret, token)
	if err != nil {
		// TODO: check log location
		log.Fatal(err)
	}
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bot/webhook" {
		http.NotFound(w, r)
		return
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.Contains(message.Text, "ルンちゃん") {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ほい")).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}

func main() {
	// initialize handler
	http.HandleFunc("/bot/webhook", webhookHandler)

	// initialize server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
