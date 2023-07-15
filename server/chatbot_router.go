package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gorilla/mux"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/bots"
	compPage "github.com/tinode/chat/server/extra/page"
	"github.com/tinode/chat/server/extra/pkg/queue"
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/ruleset/page"
	extraStore "github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/types/linkit"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/extra").Subrouter()
	s.Use(mux.CORSMethodMiddleware(r))
	// common
	s.HandleFunc("/oauth/{category}/{uid1}/{uid2}", storeOAuth)
	s.HandleFunc("/page/{id}", getPage)
	s.HandleFunc("/form", postForm).Methods(http.MethodPost)
	s.HandleFunc("/queue/stats", queueStats)
	s.HandleFunc("/p/{id}/{flag}", renderPage)
	// bot
	s.HandleFunc("/linkit", postLinkitData)

	return s
}

func newUrlRouter() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/u").Subrouter()
	s.HandleFunc("/{flag}", urlRedirect)
	return s
}

func newDownloadRouter() *mux.Router {
	dir := os.Getenv("DOWNLOAD_PATH")
	r := mux.NewRouter()
	r.PathPrefix("/d").Handler(http.StripPrefix("/d/", http.FileServer(http.Dir(dir))))
	return r
}

// handler

func storeOAuth(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	category := vars["category"]
	ui1, _ := strconv.ParseUint(vars["uid1"], 10, 64)
	ui2, _ := strconv.ParseUint(vars["uid2"], 10, 64)
	if ui1 == 0 || ui2 == 0 {
		errorResponse(rw, "path error")
		return
	}

	// code -> token
	provider := newProvider(category)
	tk, err := provider.StoreAccessToken(req)
	if err != nil {
		logs.Err.Println("router oauth", err)
		errorResponse(rw, "oauth error")
		return
	}

	// store
	extra := model.JSON{}
	_ = extra.Scan(tk["extra"])
	err = extraStore.Chatbot.OAuthSet(model.OAuth{
		UID:   types.Uid(ui1).UserId(),
		Topic: types.Uid(ui2).UserId(),
		Name:  category,
		Type:  category,
		Token: tk["token"].(string),
		Extra: extra,
	})
	if err != nil {
		logs.Err.Println("router oauth", err)
		errorResponse(rw, "store error")
		return
	}

	_, _ = rw.Write([]byte("ok"))
}

func getPage(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	p, err := extraStore.Chatbot.PageGet(id)
	if err != nil {
		logs.Err.Println(err)
		errorResponse(rw, "page error")
		return
	}

	var comp app.UI
	switch p.Type {
	case model.PageForm:
		f, _ := extraStore.Chatbot.FormGet(p.PageID)
		comp = compPage.RenderForm(p, f)
	case model.PageOkr:
		comp = compPage.RenderOkr(p)
	case model.PageTable:
		comp = compPage.RenderTable(p)
	case model.PageShare:
		comp = compPage.RenderShare(p)
	case model.PageJson:
		comp = compPage.RenderJson(p)
	case model.PageHtml:
		comp = compPage.RenderHtml(p)
	case model.PageMarkdown:
		comp = compPage.RenderMarkdown(p)
	case model.PageChart:
		d, err := json.Marshal(p.Schema)
		if err != nil {
			return
		}
		var msg extraTypes.ChartMsg
		err = json.Unmarshal(d, &msg)
		if err != nil {
			return
		}

		line := charts.NewLine()
		line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
			Title:    msg.Title,
			Subtitle: msg.SubTitle,
		}), charts.WithInitializationOpts(opts.Initialization{PageTitle: "Chart"}))

		var lineData []opts.LineData
		for _, i := range msg.Series {
			lineData = append(lineData, opts.LineData{Value: i})
		}

		line.SetXAxis(msg.XAxis).AddSeries("Chart", lineData)

		_ = line.Render(rw)
		return
	default:
		errorResponse(rw, "page error")
		return
	}

	_, _ = rw.Write([]byte(fmt.Sprintf(compPage.Layout, app.HTMLString(comp))))
}

func renderPage(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pageRuleId := vars["id"]
	flag := vars["flag"]

	p, err := extraStore.Chatbot.ParameterGet(flag)
	if err != nil {
		errorResponse(rw, "flag error")
		return
	}
	if p.IsExpired() {
		errorResponse(rw, "page expired")
		return
	}

	topic, _ := p.Params.String("topic")
	uid, _ := p.Params.String("uid")

	ctx := extraTypes.Context{
		RcptTo:     topic,
		AsUser:     types.ParseUserId(uid),
		PageRuleId: pageRuleId,
	}

	var botHandler bots.Handler
	for _, handler := range bots.List() {
		for _, item := range handler.Rules() {
			switch v := item.(type) {
			case []page.Rule:
				for _, rule := range v {
					if rule.Id == pageRuleId {
						botHandler = handler
					}
				}
			}
		}
	}

	if botHandler == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	html, err := botHandler.Page(ctx, flag)
	_, _ = rw.Write([]byte(html))
}

func postForm(rw http.ResponseWriter, req *http.Request) {
	_ = req.ParseForm()
	pf := req.PostForm

	formId := pf.Get("x-form_id")
	uid := pf.Get("x-uid")
	uid2 := pf.Get("x-topic")

	userUid := types.ParseUserId(uid)
	topicUid := types.ParseUserId(uid2)
	topic := userUid.P2PName(topicUid)

	formData, err := extraStore.Chatbot.FormGet(formId)
	if err != nil {
		return
	}
	if formData.State == model.FormStateSubmitSuccess || formData.State == model.FormStateSubmitFailed {
		return
	}

	values := make(map[string]interface{})
	d, err := json.Marshal(formData.Schema)
	if err != nil {
		return
	}
	var formMsg extraTypes.FormMsg
	err = json.Unmarshal(d, &formMsg)
	if err != nil {
		return
	}
	for _, field := range formMsg.Field {
		switch field.Type {
		case extraTypes.FormFieldCheckbox:
			value := pf[field.Key]
			switch field.ValueType {
			case extraTypes.FormFieldValueStringSlice:
				values[field.Key] = value
			case extraTypes.FormFieldValueInt64Slice:
				var tmp []int64
				for _, s := range value {
					i, _ := strconv.ParseInt(s, 10, 64)
					tmp = append(tmp, i)
				}
				values[field.Key] = tmp
			case extraTypes.FormFieldValueFloat64Slice:
				var tmp []float64
				for _, s := range value {
					i, _ := strconv.ParseFloat(s, 64)
					tmp = append(tmp, i)
				}
				values[field.Key] = tmp
			}
		default:
			value := pf.Get(field.Key)
			switch field.ValueType {
			case extraTypes.FormFieldValueString:
				values[field.Key] = value
			case extraTypes.FormFieldValueBool:
				if value == "true" {
					values[field.Key] = true
				}
				if value == "false" {
					values[field.Key] = false
				}
			case extraTypes.FormFieldValueInt64:
				values[field.Key], _ = strconv.ParseInt(value, 10, 64)
			case extraTypes.FormFieldValueFloat64:
				values[field.Key], _ = strconv.ParseFloat(value, 64)
			}
		}
	}

	ctx := extraTypes.Context{
		Original:   topicUid.UserId(),
		RcptTo:     topic,
		AsUser:     types.ParseUserId(uid),
		FormId:     formData.FormID,
		FormRuleId: formMsg.ID,
	}

	// user auth record
	_, authLvl, _, _, _ := store.Users.GetAuthRecord(userUid, "basic")

	// get bot handler
	formRuleId, ok := formData.Schema.String("id")
	if !ok {
		logs.Err.Printf("form %s %s", formId, "error form rule id")
		return
	}
	var botHandler bots.Handler
	for _, handler := range bots.List() {
		for _, item := range handler.Rules() {
			switch v := item.(type) {
			case []form.Rule:
				for _, rule := range v {
					if rule.Id == formRuleId {
						botHandler = handler
					}
				}
			}
		}
	}

	if botHandler != nil {
		if !botHandler.IsReady() {
			logs.Info.Printf("bot %s unavailable", topic)
			return
		}

		switch botHandler.AuthLevel() {
		case auth.LevelRoot:
			if authLvl != auth.LevelRoot {
				// Unauthorized
				return
			}
		}

		// form message
		payload, err := botHandler.Form(ctx, values)
		if err != nil {
			logs.Warn.Printf("topic[%s]: failed to form bot: %v", topic, err)
			return
		}

		// stats
		statsInc("BotRunFormTotal", 1)

		// send message
		if payload == nil {
			return
		}

		botSend(topic, topicUid, payload)

		// workflow form step
		workflowFlag, _ := formData.Extra.String("workflow_flag")
		workflowVersion, _ := formData.Extra.Int64("workflow_version")
		nextWorkflow(ctx, workflowFlag, int(workflowVersion), topic, topicUid)
	}

	_, _ = rw.Write([]byte("ok"))
}

func postLinkitData(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	if req.Method == http.MethodOptions {
		return
	}

	// authorization
	token := req.Header.Get("Authorization")
	token = strings.TrimSpace(token)
	token = strings.ReplaceAll(token, "Bearer ", "")

	p, err := extraStore.Chatbot.ParameterGet(token)
	if err != nil {
		errorResponse(rw, "error")
		return
	}
	if p.ID <= 0 || p.IsExpired() {
		errorResponse(rw, "401")
		return
	}

	ui1, _ := p.Params.String("uid")
	uid1 := types.ParseUserId(ui1)
	if uid1.IsZero() {
		errorResponse(rw, "401")
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		errorResponse(rw, "error")
		return
	}

	var data linkit.Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		errorResponse(rw, "error")
		return
	}

	switch data.Action {
	case linkit.Agent:
		userUid := uid1

		d, err := json.Marshal(data.Content)
		if err != nil {
			errorResponse(rw, "error")
			return
		}
		a := agent.Data{}
		err = json.Unmarshal(d, &a)
		if err != nil {
			errorResponse(rw, "error")
			return
		}

		subs, err := store.Users.FindSubs(userUid, [][]string{{"bot"}}, nil, true)
		if err != nil {
			errorResponse(rw, "error")
			return
		}

		// user auth record
		_, authLvl, _, _, _ := store.Users.GetAuthRecord(userUid, "basic")

		for _, sub := range subs {
			if !isBot(sub) {
				continue
			}

			topic := sub.User
			topicUid := types.ParseUid(topic)

			// bot name
			name := botName(sub)
			handle, ok := bots.List()[name]
			if !ok {
				continue
			}

			if !handle.IsReady() {
				logs.Info.Printf("bot %s unavailable", topic)
				continue
			}

			switch handle.AuthLevel() {
			case auth.LevelRoot:
				if authLvl != auth.LevelRoot {
					// Unauthorized
					continue
				}
			}

			ctx := extraTypes.Context{
				Original:     topicUid.UserId(),
				RcptTo:       topic,
				AsUser:       userUid,
				AgentId:      a.Id,
				AgentVersion: data.Version,
			}
			payload, err := handle.Agent(ctx, data.Content)
			if err != nil {
				logs.Warn.Printf("topic[%s]: failed to agent bot: %v", topic, err)
				continue
			}

			// stats
			statsInc("BotRunAgentTotal", 1)

			// send message
			if payload == nil {
				continue
			}

			botSend(uid1.P2PName(topicUid), topicUid, payload)
		}
	case linkit.Pull:
		list, err := extraStore.Chatbot.ListInstruct(uid1, false)
		if err != nil {
			errorResponse(rw, "error")
			return
		}
		var instruct []map[string]interface{}
		instruct = []map[string]interface{}{}
		for _, item := range list {
			instruct = append(instruct, map[string]interface{}{
				"no":        item.No,
				"bot":       item.Bot,
				"flag":      item.Flag,
				"content":   item.Content,
				"expire_at": item.ExpireAt,
			})
		}

		res, _ := json.Marshal(map[string]interface{}{
			"instruct": instruct,
		})
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write(res)
		return
	case linkit.Info:
		user, err := store.Users.Get(uid1)
		if err != nil {
			errorResponse(rw, "error")
			return
		}

		result := map[string]interface{}{
			"version":  1,
			"username": utils.Fn(user.Public),
		}
		res, _ := json.Marshal(result)
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write(res)
		return
	case linkit.Bots:
		var data []map[string]interface{}
		for name, bot := range bots.List() {
			instruct, err := bot.Instruct()
			if err != nil {
				continue
			}
			if len(instruct) <= 0 {
				continue
			}
			data = append(data, map[string]interface{}{
				"id":   name,
				"name": name,
			})
		}

		result, err := json.Marshal(map[string]interface{}{
			"bots": data,
		})
		if err != nil {
			errorResponse(rw, "error")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write(result)
		return
	case linkit.Help:
		if id, ok := data.Content.(string); ok {
			if bot, ok := bots.List()[id]; ok {
				result, err := bot.Help()
				if err != nil {
					errorResponse(rw, "error")
					return
				}
				d, _ := json.Marshal(result)
				rw.Header().Set("Content-Type", "application/json")
				_, _ = rw.Write(d)
				return
			}
			_, _ = rw.Write([]byte("{}"))
			return
		} else {
			errorResponse(rw, "error")
			return
		}
	}

	_, _ = rw.Write([]byte("ok"))
}

func urlRedirect(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	flag, ok := vars["flag"]
	if !ok {
		errorResponse(rw, "error")
		return
	}

	url, err := extraStore.Chatbot.UrlGetByFlag(flag)
	if err != nil {
		errorResponse(rw, "error")
		return
	}

	// view count
	_ = extraStore.Chatbot.UrlViewIncrease(flag)

	// redirect
	http.Redirect(rw, req, url.URL, http.StatusFound)
	return
}

func queueStats(rw http.ResponseWriter, _ *http.Request) {
	html, err := queue.Stats()
	if err != nil {
		errorResponse(rw, "queue stats error")
		return
	}
	_, _ = fmt.Fprint(rw, html)
}
