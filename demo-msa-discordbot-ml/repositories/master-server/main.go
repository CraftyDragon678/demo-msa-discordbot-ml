package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"master-server/config"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	ShardManager *shard.Manager
	Prefix       = "."
)

func init() {
	newShards := state.NewShardFunc(func(m *shard.Manager, s *state.State) {
		s.AddIntents(gateway.IntentDirectMessages)
		s.AddIntents(gateway.IntentGuildMessages)
		s.AddHandler(func(m *gateway.MessageCreateEvent) { messageCreate(s, m) })
	})

	m, err := shard.NewManager("Bot "+config.Token, newShards)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := m.Open(context.Background()); err != nil {
		log.Fatalf("Error: %v", err)
	}

	var shardNum int

	m.ForEach(func(shard shard.Shard) {
		s, ok := shard.(*state.State)
		if !ok {
			return
		}

		g, err := s.Guilds()
		if err != nil {
			return
		}

		log.Printf("Shard %d/%d (%d Guilds) Stated", shardNum, m.NumShards()-1, len(g))
	})

	ShardManager = m
}

func main() {
	defer ShardManager.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
}

func messageCreate(s *state.State, m *gateway.MessageCreateEvent) {
	if s == nil {
		return
	}

	prefix := Prefix

	if len(m.Content) <= len(prefix) ||
		m.Content[:len(prefix)] != prefix {
		return
	}

	args := strings.Fields(m.Content[len(prefix):])
	commandName := strings.ToLower(args[0])

	data, _ := json.Marshal(m)
	sendMessage(data, 50392, "/"+commandName+"?args="+strings.Join(args, ","))
	sendMessage(data, 50393, "/"+commandName+"?args="+strings.Join(args, ","))
	sendMessage(data, 50394, "/"+commandName+"?args="+strings.Join(args, ","))
}

func sendMessage(data []byte, port int, path string) {
	reader := bytes.NewReader(data)
	http.Post("http://localhost:"+strconv.Itoa(port)+path, "application/json", reader)
}
