package finance

import (
	"context"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/doctorxiong"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: `fund [string]`,
		Help:   `Get fund`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			code, _ := tokens[1].Value.String()

			reply, err := doctorxiong.GetFund(context.Background(), code)
			if err != nil {
				return nil
			}

			if reply.Name != "" {
				var xAxis []string
				var series []float64
				if reply.NetWorthDataDate == nil || len(reply.NetWorthDataDate) == 0 {
					xAxis = reply.MillionCopiesIncomeDataDate
					series = reply.MillionCopiesIncomeDataIncome
				} else {
					xAxis = reply.NetWorthDataDate
					series = reply.NetWorthDataUnit
				}

				return bots.StorePage(ctx, model.PageChart, types.ChartMsg{
					Title:    fmt.Sprintf("Fund %s (%s)", reply.Name, reply.Code),
					SubTitle: "Data for the last 90 days",
					XAxis:    xAxis,
					Series:   series,
				})
			}

			return types.TextMsg{Text: "failed"}
		},
	},
	{
		Define: `stock [string]`,
		Help:   `Get stock`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			code, _ := tokens[1].Value.String()

			reply, err := doctorxiong.GetStock(context.Background(), code)
			if err != nil {
				return nil
			}

			return types.InfoMsg{
				Title: fmt.Sprintf("Stock %s", code),
				Model: reply,
			}
		},
	},
}
