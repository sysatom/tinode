package command

import (
	"context"
	"fmt"
	"github.com/tinode/chat/server/extra/types"
	"strings"
	"unicode"
)

type Rule struct {
	Define  string
	Help    string
	Handler func(context.Context, []*Token) []types.MsgPayload
}

type Ruleset []Rule

func (r Ruleset) Help(in string) ([]types.MsgPayload, error) {
	if strings.ToLower(in) == "help" {
		var helpMsg string
		for _, rule := range r {
			helpMsg = fmt.Sprintf("%s%s%s%s\n", helpMsg, rule.Define, " :: ", rule.Help)
		}
		return []types.MsgPayload{
			types.TextMsg{Text: strings.TrimLeftFunc(helpMsg, unicode.IsSpace)},
		}, nil
	}
	return nil, nil
}

func (r Ruleset) ProcessCommand(ctx context.Context, in string) ([]types.MsgPayload, error) {
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
