package router

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store/types"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var rOAuth = regexp.MustCompile(`/oauth/\w+`)
var rOAuthRedirect = regexp.MustCompile(`/oauth/\w+/redirect`)
var rChart = regexp.MustCompile(`/chart/\w+`)

func ServeExtra(rw http.ResponseWriter, req *http.Request) {
	switch {
	case rOAuth.MatchString(req.URL.Path):
		oauth(rw, req)
	case rOAuthRedirect.MatchString(req.URL.Path):
		oauthRedirect(rw, req)
	case rChart.MatchString(req.URL.Path):
		chart(rw, req)
	default:
		rw.Write([]byte("Unknown Pattern"))
	}
}

// handler

func oauthRedirect(rw http.ResponseWriter, req *http.Request) {
	category := strings.ReplaceAll(req.URL.Path, "/extra/oauth/", "")
	category = strings.ReplaceAll(req.URL.Path, "/redirect", "")
	category = strings.ToLower(category)
	provider := newProvider(category)
	url, err := provider.Redirect(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("oauth redirect error"))
		return
	}
	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusFound)
}

func oauth(rw http.ResponseWriter, req *http.Request) {
	paramsPatch := strings.ToLower(strings.ReplaceAll(req.URL.Path, "/extra/oauth/", ""))
	params := strings.Split(paramsPatch, "/")
	if len(params) != 3 {
		rw.Write([]byte("path error"))
		return
	}
	ui1, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		logs.Err.Println("router oauth", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("path error"))
		return
	}
	ui2, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		logs.Err.Println("router oauth", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("path error"))
		return
	}

	// code -> token
	provider := newProvider(params[0])
	tk, err := provider.StoreAccessToken(req)
	if err != nil {
		logs.Err.Println("router oauth", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("oauth error"))
		return
	}

	// store
	extra := model.JSON{}
	_ = extra.Scan(tk["extra"])
	err = store.Chatbot.OAuthSet(model.OAuth{
		Uid:   types.Uid(ui1).UserId(),
		Topic: types.Uid(ui2).UserId(),
		Name:  params[0],
		Type:  params[0],
		Token: tk["token"].(string),
		Extra: extra,
	})
	if err != nil {
		logs.Err.Println("router oauth", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("store error"))
		return
	}

	rw.Write([]byte("ok"))
}

func chart(rw http.ResponseWriter, _ *http.Request) {
	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "My first bar chart generated by go-echarts",
		Subtitle: "It's extremely easy to use, right?",
	}), charts.WithInitializationOpts(opts.Initialization{PageTitle: "Chart"}))

	// Put data into instance
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems())

	bar.Render(rw)
}
