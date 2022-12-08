package help

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/nikolaydubina/calendarheatmap/charts"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"image/color"
	"math/big"
	"strconv"
	"time"
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
		Define: "rand [number] [number]",
		Help:   `Generate random numbers`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			min, _ := tokens[1].Value.Int64()
			max, _ := tokens[2].Value.Int64()

			nBing, err := rand.Int(rand.Reader, big.NewInt(max+1-min))
			if err != nil {
				logs.Err.Println("bot command rand [number] [number]", err)
				return nil
			}
			t := nBing.Int64() + min

			return types.TextMsg{Text: strconv.FormatInt(t, 10)}
		},
	},
	{
		Define: "id",
		Help:   `Generate random id`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: types.Id().String()}
		},
	},
	{
		Define: "uid [string]",
		Help:   `Decode UID string`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			str, _ := tokens[1].Value.String()
			var uid serverTypes.Uid
			var result string
			err := uid.UnmarshalText([]byte(str))
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("%d", uid)
			}

			return types.TextMsg{Text: result}
		},
	},
	{
		Define: "ts [number]",
		Help:   `timestamp format`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			num, _ := tokens[1].Value.Int64()
			t := time.Unix(num, 0)
			return types.TextMsg{Text: t.Format(time.RFC3339)}
		},
	},
	{
		Define: "messages",
		Help:   `Demo messages`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: `qr [string]`,
		Help:   `Generate QR code`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: `pinyin [string]`,
		Help:   "chinese pinyin conversion",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: "form",
		Help:   `Demo form`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.StoreForm(ctx, types.FormMsg{
				ID:    helpFormID,
				Title: "Current Value: 1, add/reduce ?",
				Field: []types.FormField{
					{
						Key:         "text",
						Type:        types.FormFieldText,
						ValueType:   types.FormFieldValueString,
						Value:       "",
						Label:       "Text",
						Placeholder: "Input text",
					},
					{
						Key:         "password",
						Type:        types.FormFieldPassword,
						ValueType:   types.FormFieldValueString,
						Value:       "",
						Label:       "Password",
						Placeholder: "Input password",
						Required:    true,
					},
					{
						Key:         "number",
						Type:        types.FormFieldNumber,
						ValueType:   types.FormFieldValueInt64,
						Value:       "",
						Label:       "Number",
						Placeholder: "Input number",
					},
					{
						Key:         "bool",
						Type:        types.FormFieldRadio,
						ValueType:   types.FormFieldValueBool,
						Value:       "",
						Label:       "Bool",
						Placeholder: "Switch",
						Option:      []string{"true", "false"},
					},
					{
						Key:         "multi",
						Type:        types.FormFieldCheckbox,
						ValueType:   types.FormFieldValueStringSlice,
						Value:       "",
						Label:       "Multiple",
						Placeholder: "Select multiple",
						Option:      []string{"a", "b", "c"},
					},
					{
						Key:         "textarea",
						Type:        types.FormFieldTextarea,
						ValueType:   types.FormFieldValueString,
						Value:       "",
						Label:       "Textarea",
						Placeholder: "Input textarea",
					},
					{
						Key:         "select",
						Type:        types.FormFieldSelect,
						ValueType:   types.FormFieldValueFloat64,
						Value:       "",
						Label:       "Select",
						Placeholder: "Select float",
						Option:      []string{"1.01", "2.02", "3.03"},
					},
					{
						Key:         "range",
						Type:        types.FormFieldRange,
						ValueType:   types.FormFieldValueInt64,
						Value:       "",
						Label:       "Range",
						Placeholder: "range value",
					},
				},
			})
		},
	},
	{
		Define: "action",
		Help:   "Demo action",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.ActionMsg{
				ID:     helpActionID,
				Title:  "Operate ... ?",
				Option: []string{"do1", "do2"},
				Value:  "",
			}
		},
	},
	{
		Define: "echo [any]",
		Help:   "print",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			val := tokens[1].Value.Source
			return types.TextMsg{Text: fmt.Sprintf("%v", val)}
		},
	},
	{
		Define: "agent",
		Help:   `agent url`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.AgentURI(ctx)
		},
	},
	{
		Define: "plot",
		Help:   `plot graph`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			p := plot.New()

			p.Title.Text = "Plotutil example"
			p.X.Label.Text = "X"
			p.Y.Label.Text = "Y"

			err := plotutil.AddLinePoints(p,
				"First", randomPoints(15),
				"Second", randomPoints(15),
				"Third", randomPoints(15))
			if err != nil {
				panic(err)
			}

			w := bytes.NewBufferString("")

			c := vgimg.New(vg.Points(500), vg.Points(500))
			dc := draw.New(c)
			p.Draw(dc)

			png := vgimg.PngCanvas{Canvas: c}
			if _, err := png.WriteTo(w); err != nil {
				panic(err)
			}

			raw := base64.StdEncoding.EncodeToString(w.Bytes())

			return types.ImageMsg{
				Width:       500,
				Height:      500,
				Alt:         "Plot.png",
				Mime:        "image/png",
				Size:        w.Len(),
				ImageBase64: raw,
			}
		},
	},
	{
		Define: "heatmap",
		Help:   `heatmap`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			fontFace, err := charts.LoadFontFace(defaultFontFaceBytes, opentype.FaceOptions{
				Size:    26,
				DPI:     280,
				Hinting: font.HintingNone,
			})
			if err != nil {
				logs.Err.Println(err)
				return nil
			}

			var colorscale charts.BasicColorScale
			colorscale, err = charts.NewBasicColorscaleFromCSV(bytes.NewBuffer(defaultColorScaleBytes))
			if err != nil {
				logs.Err.Println(err)
				return nil
			}

			// data
			counts := map[string]int{
				"2020-05-16": 8,
				"2020-05-17": 13,
				"2020-05-18": 5,
				"2020-05-19": 8,
				"2020-05-20": 5,
			}

			conf := charts.HeatmapConfig{
				Counts:              counts,
				ColorScale:          colorscale,
				DrawMonthSeparator:  true,
				DrawLabels:          true,
				Margin:              30,
				BoxSize:             150,
				MonthSeparatorWidth: 5,
				MonthLabelYOffset:   50,
				TextWidthLeft:       300,
				TextHeightTop:       200,
				TextColor:           color.RGBA{100, 100, 100, 255},
				BorderColor:         color.RGBA{200, 200, 200, 255},
				Locale:              "en_US",
				Format:              "png",
				FontFace:            fontFace,
				ShowWeekdays: map[time.Weekday]bool{
					time.Monday:    true,
					time.Wednesday: true,
					time.Friday:    true,
				},
			}
			w := bytes.NewBufferString("")
			_ = charts.WriteHeatmap(conf, w)

			raw := base64.StdEncoding.EncodeToString(w.Bytes())

			return types.ImageMsg{
				Width:       1858,
				Height:      275,
				Alt:         "Heatmap.png",
				Mime:        "image/png",
				Size:        w.Len(),
				ImageBase64: raw,
			}
		},
	},
}

//go:embed fonts/Sunflower-Medium.ttf
var defaultFontFaceBytes []byte

//go:embed colorscales/green-blue-9.csv
var defaultColorScaleBytes []byte

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		num, _ := rand.Int(rand.Reader, big.NewInt(100))
		if i == 0 {
			pts[i].X = float64(num.Int64())
		} else {
			pts[i].X = pts[i-1].X + float64(num.Int64())
		}
		pts[i].Y = pts[i].X + 10*float64(num.Int64())
	}
	return pts
}
