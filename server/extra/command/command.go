package command

import (
	"github.com/tinode/chat/server/extra/types"
	"strings"
)

type Rule struct {
	Define  string
	Help    string
	Handler func(types.Context, []*Token) []types.MsgPayload
}

type Ruleset []Rule

func (r Ruleset) Help(in string) ([]types.MsgPayload, error) {
	if strings.ToLower(in) == "help" {
		table := types.TableMsg{
			Header: []string{"Define", "Help"},
		}
		for _, rule := range r {
			table.Row = append(table.Row, []interface{}{rule.Define, rule.Help})
		}
		return []types.MsgPayload{table}, nil
	}
	return nil, nil
}

func (r Ruleset) ProcessCommand(ctx types.Context, in string) ([]types.MsgPayload, error) {
	var result []types.MsgPayload
	for _, rule := range r {
		tokens, err := ParseCommand(in)
		if err != nil {
			return nil, err
		}
		check, err := SyntaxCheck(rule.Define, tokens)
		if err != nil {
			return nil, err
		}
		if !check {
			continue
		}

		if ret := rule.Handler(ctx, tokens); len(ret) > 0 {
			result = ret
		}
	}
	return result, nil
}
