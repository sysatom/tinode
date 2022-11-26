package mtg

import (
	"context"
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/scryfall"
	"github.com/tinode/chat/server/logs"
)

var commandRules = []command.Rule{
	{
		Define: "info",
		Help:   `Bot info`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return nil
		},
	},
	{
		Define: "search [string]",
		Help:   `Search cards.`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			keyword, _ := tokens[1].Value.String()
			provider := scryfall.NewScryfall()
			result, err := provider.CardsSearch(context.Background(), fmt.Sprintf("%s lang:zhs", keyword))
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "search error"}
			}
			if len(result) == 0 {
				return types.TextMsg{Text: "empty"}
			}
			var cards []types.CardMsg
			for _, card := range result {
				name := card.PrintedName
				if name == "" {
					name = card.Name
				}
				cards = append(cards, types.CardMsg{
					Name: name,
					URI:  card.ScryfallURI,
				})
			}
			return types.CardListMsg{
				Cards: cards,
			}
		},
	},
}
