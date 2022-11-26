package main

import (
	"github.com/elastic/go-sysinfo"
	"github.com/robfig/cron/v3"
	"github.com/tinode/chat/server/extra/agent"
	"github.com/tinode/chat/server/extra/bots/server"
	"log"
)

var agentURI string

func main() {
	// args
	agentURI = agent.URI()

	// cron
	c := cron.New()
	_, err := c.AddFunc("* * * * *", stats)
	if err != nil {
		panic(err)
	}
	c.Run()
}

func stats() {
	host, err := sysinfo.Host()
	if err != nil {
		log.Println(err)
		return
	}
	info := host.Info()
	cpu, _ := host.CPUTime()
	memory, _ := host.Memory()

	err = agent.PostData(agentURI, agent.Data{
		Id:      server.StatsAgentID,
		Version: server.AgentVersion,
		Content: map[string]interface{}{
			"info":   info,
			"cpu":    cpu,
			"memory": memory,
		},
	})
	if err != nil {
		log.Println(err)
	}
}
