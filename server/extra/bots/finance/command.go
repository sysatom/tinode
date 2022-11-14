package finance

import (
	"context"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/doctorxiong"
	"net/http"
	"strconv"
	"time"
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

			reply, err := getFund(context.Background(), code)
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

			reply, err := getStock(context.Background(), code)
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

func getFund(ctx context.Context, code string) (*FundReply, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -90).Format("2006-01-02")
	endDate := now.Format("2006-01-02")
	dx := doctorxiong.NewDoctorxiong("")
	resp, err := dx.GetFundDetail(ctx, code, startDate, endDate)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	fund := resp.Data
	var netWorthDataDate []string
	var netWorthDataUnit []float64
	var netWorthDataIncrease []float64
	for _, item := range fund.NetWorthData {
		netWorthDataDate = append(netWorthDataDate, item[0].(string))
		f1, _ := strconv.ParseFloat(item[1].(string), 64)
		netWorthDataUnit = append(netWorthDataUnit, f1)
		f2, _ := strconv.ParseFloat(item[2].(string), 64)
		netWorthDataIncrease = append(netWorthDataIncrease, f2)
	}
	var millionCopiesIncomeDataDate []string
	var millionCopiesIncomeDataIncome []float64
	for _, item := range fund.MillionCopiesIncomeData {
		millionCopiesIncomeDataDate = append(millionCopiesIncomeDataDate, item[0].(string))
		f1, _ := strconv.ParseFloat(item[1].(string), 64)
		millionCopiesIncomeDataIncome = append(millionCopiesIncomeDataIncome, f1)
	}

	return &FundReply{
		Code:                          fund.Code,
		Name:                          fund.Name,
		Type:                          fund.Type,
		NetWorth:                      fund.NetWorth,
		ExpectWorth:                   fund.ExpectWorth,
		TotalWorth:                    fund.TotalWorth,
		ExpectGrowth:                  fund.ExpectGrowth,
		DayGrowth:                     fund.DayGrowth,
		LastWeekGrowth:                fund.LastWeekGrowth,
		LastMonthGrowth:               fund.LastMonthGrowth,
		LastThreeMonthsGrowth:         fund.LastThreeMonthsGrowth,
		LastSixMonthsGrowth:           fund.LastSixMonthsGrowth,
		LastYearGrowth:                fund.LastYearGrowth,
		BuyMin:                        fund.BuyMin,
		BuySourceRate:                 fund.BuySourceRate,
		BuyRate:                       fund.BuyRate,
		Manager:                       fund.Manager,
		FundScale:                     fund.FundScale,
		NetWorthDate:                  fund.NetWorthDate,
		ExpectWorthDate:               fund.ExpectWorthDate,
		NetWorthDataDate:              netWorthDataDate,
		NetWorthDataUnit:              netWorthDataUnit,
		NetWorthDataIncrease:          netWorthDataIncrease,
		MillionCopiesIncomeDate:       "",
		SevenDaysYearIncome:           0,
		MillionCopiesIncomeDataDate:   millionCopiesIncomeDataDate,
		MillionCopiesIncomeDataIncome: millionCopiesIncomeDataIncome,
	}, nil
}

func getStock(ctx context.Context, code string) (*StockReply, error) {
	dx := doctorxiong.NewDoctorxiong("")
	resp, err := dx.GetStockDetail(ctx, code)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	stock := resp.Data
	if len(resp.Data) <= 0 {
		return &StockReply{}, nil
	}
	return &StockReply{
		Code:             stock[0].Code,
		Name:             stock[0].Name,
		Type:             stock[0].Type,
		PriceChange:      stock[0].PriceChange,
		ChangePercent:    stock[0].ChangePercent,
		Open:             stock[0].Open,
		Close:            stock[0].Close,
		Price:            stock[0].Price,
		High:             stock[0].High,
		Low:              stock[0].Low,
		Volume:           stock[0].Volume,
		Turnover:         stock[0].Turnover,
		TurnoverRate:     stock[0].TurnoverRate,
		TotalWorth:       stock[0].TotalWorth,
		CirculationWorth: stock[0].CirculationWorth,
		Date:             stock[0].Date,
		Buy:              nil,
		Sell:             nil,
		Pb:               stock[0].Pb,
		Spe:              stock[0].Spe,
		Pe:               stock[0].Pe,
	}, nil
}

type FundReply struct {
	Code                          string    `json:"code,omitempty"`
	Name                          string    `json:"name,omitempty"`
	Type                          string    `json:"type,omitempty"`
	NetWorth                      float64   `json:"net_worth,omitempty"`
	ExpectWorth                   float64   `json:"expect_worth,omitempty"`
	TotalWorth                    float64   `json:"total_worth,omitempty"`
	ExpectGrowth                  string    `json:"expect_growth,omitempty"`
	DayGrowth                     string    `json:"day_growth,omitempty"`
	LastWeekGrowth                string    `json:"last_week_growth,omitempty"`
	LastMonthGrowth               string    `json:"last_month_growth,omitempty"`
	LastThreeMonthsGrowth         string    `json:"last_three_months_growth,omitempty"`
	LastSixMonthsGrowth           string    `json:"last_six_months_growth,omitempty"`
	LastYearGrowth                string    `json:"last_year_growth,omitempty"`
	BuyMin                        string    `json:"buy_min,omitempty"`
	BuySourceRate                 string    `json:"buy_source_rate,omitempty"`
	BuyRate                       string    `json:"buy_rate,omitempty"`
	Manager                       string    `json:"manager,omitempty"`
	FundScale                     string    `json:"fund_scale,omitempty"`
	NetWorthDate                  string    `json:"net_worth_date,omitempty"`
	ExpectWorthDate               string    `json:"expect_worth_date,omitempty"`
	NetWorthDataDate              []string  `json:"net_worth_data_date,omitempty"`
	NetWorthDataUnit              []float64 `json:"net_worth_data_unit,omitempty"`
	NetWorthDataIncrease          []float64 `json:"net_worth_data_increase,omitempty"`
	MillionCopiesIncomeDataDate   []string  `json:"million_copies_income_data_date,omitempty"`
	MillionCopiesIncomeDataIncome []float64 `json:"million_copies_income_data_income,omitempty"`
	MillionCopiesIncomeDate       string    `json:"million_copies_income_date,omitempty"`
	SevenDaysYearIncome           float64   `json:"seven_days_year_income,omitempty"`
}

type StockReply struct {
	Code             string   `json:"code,omitempty"`
	Name             string   `json:"name,omitempty"`
	Type             string   `json:"type,omitempty"`
	PriceChange      string   `json:"price_change,omitempty"`
	ChangePercent    string   `json:"change_percent,omitempty"`
	Open             string   `json:"open,omitempty"`
	Close            string   `json:"close,omitempty"`
	Price            string   `json:"price,omitempty"`
	High             string   `json:"high,omitempty"`
	Low              string   `json:"low,omitempty"`
	Volume           string   `json:"volume,omitempty"`
	Turnover         string   `json:"turnover,omitempty"`
	TurnoverRate     string   `json:"turnover_rate,omitempty"`
	TotalWorth       string   `json:"total_worth,omitempty"`
	CirculationWorth string   `json:"circulation_worth,omitempty"`
	Date             string   `json:"date,omitempty"`
	Buy              []string `json:"buy,omitempty"`
	Sell             []string `json:"sell,omitempty"`
	Pb               string   `json:"pb,omitempty"`
	Spe              string   `json:"spe,omitempty"`
	Pe               string   `json:"pe,omitempty"`
}
