package main

import (
	"github.com/robfig/cron/v3"
	"github.com/tinode/chat/server/extra/agent"
	"github.com/tinode/chat/server/extra/bots/clipboard"
	"log"
)

var agentURI string

func main() {
	agent.StartInfo()

	// args
	agentURI = agent.URI()

	// cron
	c := cron.New()
	_, err := c.AddFunc("* * * * *", example)
	if err != nil {
		panic(err)
	}
	c.Run()
}

func example() {
	err := agent.PostData(agentURI, agent.Data{
		Id:      clipboard.UploadAgentID,
		Content: nil,
	})
	if err != nil {
		log.Println(err)
	}
}
