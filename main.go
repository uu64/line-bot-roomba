package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

type iftttReqBody struct {
	Event string `json:event`
}

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

func iftttHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bot/ifttt" {
		http.NotFound(w, r)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var body iftttReqBody
	if err = json.Unmarshal(b, &body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if body.Event == "start-cleaning" {
		if _, err := bot.BroadcastMessage(linebot.NewTextMessage("掃除おわった")).Do(); err != nil {
			log.Print(err)
		}
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
	http.HandleFunc("/bot/ifttt", iftttHandler)

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
