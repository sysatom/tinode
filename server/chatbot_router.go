package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gorilla/mux"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/page"
	"github.com/tinode/chat/server/extra/pkg/queue"
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/types/helper"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
	serverStore "github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"io"
	"net/http"
	"os"
	"strconv"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/extra").Subrouter()
	s.Use(mux.CORSMethodMiddleware(r))
	s.HandleFunc("/oauth/{category}/{uid1}/{uid2}", storeOAuth)
	s.HandleFunc("/page/{id}", getPage)
	s.HandleFunc("/form", postForm).Methods(http.MethodPost)
	s.HandleFunc("/webhook/{uid1}/{uid2}/{uid3}", webhook).Methods(http.MethodPost)
	s.HandleFunc("/helper/{uid1}/{uid2}", postHelper).Methods(http.MethodPost)
	s.HandleFunc("/queue/stats", queueStats)

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
		errorResponse(rw, "store error")
		return
	}

	_, _ = rw.Write([]byte("ok"))
}

func getPage(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	p, err := store.Chatbot.PageGet(id)
	if err != nil {
		logs.Err.Println(err)
		errorResponse(rw, "page error")
		return
	}

	var comp app.UI
	switch p.Type {
	case model.PageForm:
		f, _ := store.Chatbot.FormGet(p.PageId)
		comp = page.RenderForm(p, f)
	case model.PageOkr:
		comp = page.RenderOkr(p)
	case model.PageTable:
		comp = page.RenderTable(p)
	case model.PageShare:
		comp = page.RenderShare(p)
	case model.PageJson:
		comp = page.RenderJson(p)
	case model.PageHtml:
		comp = page.RenderHtml(p)
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

	_, _ = rw.Write([]byte(fmt.Sprintf(page.Layout, app.HTMLString(comp))))
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

	formData, err := store.Chatbot.FormGet(formId)
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
		FormId:     formData.FormId,
		FormRuleId: formMsg.ID,
	}

	// user auth record
	_, authLvl, _, _, _ := serverStore.Users.GetAuthRecord(userUid, "basic")

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
	}

	_, _ = rw.Write([]byte("ok"))
}

func webhook(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	ui1, _ := strconv.ParseUint(vars["uid1"], 10, 64)
	ui2, _ := strconv.ParseUint(vars["uid2"], 10, 64)
	ui3, _ := strconv.ParseUint(vars["uid3"], 10, 64)

	uid1 := types.Uid(ui1)
	uid2 := types.Uid(ui2)
	uid3 := types.Uid(ui3)
	topic := uid1.P2PName(uid2)

	value, err := store.Chatbot.DataGet(uid1, uid2.UserId(), fmt.Sprintf("webhook:%s", uid3.String()))
	if err != nil {
		errorResponse(rw, "webhook error")
		return
	}
	_, ok := value.String("value")
	if !ok {
		errorResponse(rw, "webhook error")
		return
	}

	d, _ := io.ReadAll(req.Body)

	txt := ""
	if len(d) > 1000 {
		txt = fmt.Sprintf("[webhook:%s] body too long", uid3.String())
	} else {
		txt = fmt.Sprintf("[webhook:%s] %s", uid3.String(), string(d))
	}
	botSend(topic, uid2, extraTypes.TextMsg{Text: txt})
	_, _ = rw.Write([]byte("ok"))
}

func postHelper(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(req)

	ui1, _ := strconv.ParseUint(vars["uid1"], 10, 64)
	ui2, _ := vars["uid2"]

	uid1 := types.Uid(ui1)

	value, err := store.Chatbot.ConfigGet(uid1, "", fmt.Sprintf("helper:%d", ui1))
	if err != nil {
		errorResponse(rw, "error")
		return
	}
	uiValue, ok := value.String("value")
	if !ok {
		errorResponse(rw, "error")
		return
	}

	if uiValue != ui2 {
		errorResponse(rw, "auth error")
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		errorResponse(rw, "error")
		return
	}

	var data helper.Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		errorResponse(rw, "error")
		return
	}

	switch data.Action {
	case helper.Agent:
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

		subs, err := serverStore.Users.FindSubs(userUid, [][]string{{"bot"}}, nil, true)
		if err != nil {
			errorResponse(rw, "error")
			return
		}

		// user auth record
		_, authLvl, _, _, _ := serverStore.Users.GetAuthRecord(userUid, "basic")

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
	case helper.Pull:
		list, err := store.Chatbot.ListInstruct(uid1, false)
		if err != nil {
			errorResponse(rw, "error")
			return
		}
		var instruct = []map[string]interface{}{}
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
		_, _ = rw.Write(res)
		return
	case helper.Info:
		user, err := serverStore.Users.Get(uid1)
		if err != nil {
			errorResponse(rw, "error")
			return
		}

		result := map[string]interface{}{
			"version":  1,
			"username": utils.Fn(user.Public),
		}
		res, _ := json.Marshal(result)
		_, _ = rw.Write(res)
		return
	case helper.Bots:
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

		_, _ = rw.Write(result)
		return
	case helper.Help:
		if id, ok := data.Content.(string); ok {
			if bot, ok := bots.List()[id]; ok {
				result, err := bot.Help()
				if err != nil {
					errorResponse(rw, "error")
					return
				}
				d, _ := json.Marshal(result)
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

	url, err := store.Chatbot.UrlGetByFlag(flag)
	if err != nil {
		errorResponse(rw, "error")
		return
	}

	// view count
	_ = store.Chatbot.UrlViewIncrease(flag)

	// redirect
	http.Redirect(rw, req, url.Url, http.StatusFound)
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
