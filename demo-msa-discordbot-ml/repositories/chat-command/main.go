package main

import (
	"chat-command/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
)

var (
	client *api.Client
)

type message struct {
	Message string `json:"message"`
}

func init() {
	client = api.NewClient("Bot " + config.Token)
}

func main() {
	http.ListenAndServe(":50394", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Split(r.URL.Path, "/")[1] != "chat" {
			w.Write([]byte("Hello, World!"))
			return
		}
		res, _ := ioutil.ReadAll(r.Body)
		var m gateway.MessageCreateEvent
		json.Unmarshal(res, &m)

		query := url.Values{"message": []string{strings.Join(strings.Split(r.URL.Query().Get("args"), ","), " ")}}
		chatRes, err := http.Get("http://localhost:8000/chat?" + query.Encode())
		if err != nil {
			panic("펑")
		}
		defer chatRes.Body.Close()
		bytes, _ := ioutil.ReadAll(chatRes.Body)
		var reply message
		json.Unmarshal(bytes, &reply)
		client.SendMessage(m.ChannelID, reply.Message)

		w.Write([]byte("Hello, World!"))
	}))
}
