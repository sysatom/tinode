package dev

import "github.com/tinode/chat/server/extra/ruleset/instruct"

const (
	ExampleInstructID = "dev_example"
)

var instructRules = []instruct.Rule{
	{
		Id:   ExampleInstructID,
		Args: []string{"txt"},
	},
}
