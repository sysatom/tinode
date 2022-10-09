package command

import (
	"context"
	"github.com/tinode/chat/server/extra/types"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexRule(t *testing.T) {
	testRules := []Rule{
		{
			Define: `test`,
			Help:   `Test info`,
			Handler: func(ctx context.Context, tokens []*Token) []types.MsgPayload {
				return []types.MsgPayload{types.TextMsg{Text: "test"}}
			},
		},
		{
			Define: `todo [string]`,
			Help:   `todo something`,
			Handler: func(ctx context.Context, tokens []*Token) []types.MsgPayload {
				text, _ := tokens[1].Value.String()
				return []types.MsgPayload{types.TextMsg{Text: text}}
			},
		},
		{
			Define: `add [number] [number]`,
			Help:   `Addition`,
			Handler: func(ctx context.Context, tokens []*Token) []types.MsgPayload {
				tt1, _ := tokens[1].Value.Int64()
				tt2, _ := tokens[2].Value.Int64()
				return []types.MsgPayload{types.TextMsg{Text: strconv.Itoa(int(tt1 + tt2))}}
			},
		},
	}

	b := Ruleset(testRules)

	out, err := b.ProcessCommand(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	require.Contains(t, out, types.TextMsg{Text: "test"})

	out2, err := b.ProcessCommand(context.Background(), "add 1 2")
	if err != nil {
		t.Fatal(err)
	}
	require.Contains(t, out2, types.TextMsg{Text: "3"})

	out3, err := b.ProcessCommand(context.Background(), "help")
	if err != nil {
		t.Fatal(err)
	}
	require.Len(t, out3, 0)

	help, err := b.Help("help")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(help) > 0)

	out4, err := b.ProcessCommand(context.Background(), `todo "a b c"`)
	if err != nil {
		t.Fatal(err)
	}
	require.Contains(t, out4, types.TextMsg{Text: "a b c"})
}
