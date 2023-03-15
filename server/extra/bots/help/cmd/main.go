package main

import (
	"github.com/robfig/cron/v3"
	"github.com/tinode/chat/server/extra/agent"
	"github.com/tinode/chat/server/extra/bots/help"
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
		Id: help.ImportAgentID,
		//Version: help.AgentVersion,
		Content: nil,
	})
	if err != nil {
		log.Println(err)
	}
}
