package iot

import (
    "github.com/tinode/chat/server/extra/ruleset/agent"
    "github.com/tinode/chat/server/extra/types"
)

const (
    AgentVersion  = 1
    ExampleAgentID = "iot_example_agent"
)

var agentRules = []agent.Rule{
    {
        Id: ExampleAgentID,
        Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
            return nil
        },
    },
}
