package router

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gorilla/mux"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/router/page"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store/types"
	"net/http"
	"strconv"
	"strings"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/extra/oauth/{key}/redirect", oauthRedirect)
	r.HandleFunc("/extra/oauth/{category}/{uid1}/{uid2}", oauth)
	r.HandleFunc("/extra/chart/{key}", chart)
	r.HandleFunc("/extra/form", form)
	return r
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
	vars := mux.Vars(req)
	category := vars["category"]
	ui1, err := strconv.ParseUint(vars["uid1"], 10, 64)
	if err != nil {
		logs.Err.Println("router oauth", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("path error"))
		return
	}
	ui2, err := strconv.ParseUint(vars["uid2"], 10, 64)
	if err != nil {
		logs.Err.Println("router oauth", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("path error"))
		return
	}

	// code -> token
	provider := newProvider(category)
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
		Name:  category,
		Type:  category,
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

func form(rw http.ResponseWriter, _ *http.Request) {
	html := `
<!DOCTYPE html>
<html>
    <head>
        <title>Page</title>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
     	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.15.12/dist/css/uikit.min.css" />
		<script src="https://cdn.jsdelivr.net/npm/uikit@3.15.12/dist/js/uikit.min.js"></script>
		<script src="https://cdn.jsdelivr.net/npm/uikit@3.15.12/dist/js/uikit-icons.min.js"></script>
    </head>

    <body>
        <div id="app" style="padding: 20px">%s</div>
    </body>
</html>
`
	rw.Write([]byte(fmt.Sprintf(html, app.HTMLString(&page.Form{}))))
}
