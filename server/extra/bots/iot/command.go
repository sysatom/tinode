package iot

import (
    "github.com/tinode/chat/server/extra/ruleset/command"
    "github.com/tinode/chat/server/extra/types"
)

var commandRules = []command.Rule{
    {
        Define: "info",
        Help:   `Bot info`,
        Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
            return nil
        },
    },
}
