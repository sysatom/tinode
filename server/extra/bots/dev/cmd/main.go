package main

import (
	"github.com/robfig/cron/v3"
	"github.com/tinode/chat/server/extra/agent"
	"github.com/tinode/chat/server/extra/bots/dev"
	"log"
)

var agentURI string

func main() {
	agent.StartInfo()

	// args
	agentURI = agent.URI()

	// cron
	c := cron.New()
	_, err := c.AddFunc("* * * * *", demo)
	if err != nil {
		panic(err)
	}
	c.Run()
}

func demo() {
	err := agent.PostData(agentURI, agent.Data{
		Id: dev.ImportAgentID,
		//Version: dev.AgentVersion,
		Content: nil,
	})
	if err != nil {
		log.Println(err)
	}
}
