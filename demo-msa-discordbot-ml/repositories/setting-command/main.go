package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"setting-command/config"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
)

var (
	client *api.Client
)

func init() {
	client = api.NewClient("Bot " + config.Token)
}

func main() {
	http.ListenAndServe(":50393", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Split(r.URL.Path, "/")[1] != "setting" {
			w.Write([]byte("Hello, World!"))
			return
		}
		res, _ := ioutil.ReadAll(r.Body)
		var m gateway.MessageCreateEvent
		json.Unmarshal(res, &m)
		client.SendMessage(m.ChannelID, "Setting command")

		w.Write([]byte("Hello, World!"))
	}))
}
