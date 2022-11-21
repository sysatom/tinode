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
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverStore "github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"io"
	"net/http"
	"strconv"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/extra").Subrouter()
	s.HandleFunc("/oauth/{category}/{uid1}/{uid2}", storeOAuth)
	s.HandleFunc("/page/{id}", getPage)
	s.HandleFunc("/form", postForm).Methods(http.MethodPost)
	s.HandleFunc("/webhook/{uid1}/{uid2}/{uid3}", webhook).Methods(http.MethodPost)
	return s
}

// handler

func storeOAuth(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	category := vars["category"]
	ui1, _ := strconv.ParseUint(vars["uid1"], 10, 64)
	ui2, _ := strconv.ParseUint(vars["uid2"], 10, 64)
	if ui1 == 0 || ui2 == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("path error"))
		return
	}

	// code -> token
	provider := newProvider(category)
	tk, err := provider.StoreAccessToken(req)
	if err != nil {
		logs.Err.Println("router oauth", err)
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("oauth error"))
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
		_, _ = rw.Write([]byte("store error"))
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
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("page error"))
		return
	}

	var comp app.UI
	switch p.Type {
	case model.PageForm:
		comp = page.RenderForm(p)
	case model.PageOkr:
		comp = page.RenderOkr(p)
	case model.PageTable:
		comp = page.RenderTable(p)
	case model.PageShare:
		comp = page.RenderShare(p)
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
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("page error"))
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

	form, err := store.Chatbot.FormGet(formId)
	if err != nil {
		return
	}
	if form.State == model.FormStateSubmitSuccess || form.State == model.FormStateSubmitFailed {
		return
	}

	subs, err := serverStore.Topics.GetUsers(topic, nil)
	if err != nil {
		logs.Err.Println("hook bot incoming", err)
		return
	}

	values := make(map[string]interface{})

	d, err := json.Marshal(form.Schema)
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
		FormId:     form.FormId,
		FormRuleId: formMsg.ID,
	}

	// user auth record
	_, authLvl, _, _, _ := serverStore.Users.GetAuthRecord(userUid, "basic")

	for _, sub := range subs {
		if !isBot(sub) {
			continue
		}

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

		// form message
		payload, err := handle.Form(ctx, values)
		if err != nil {
			logs.Warn.Printf("topic[%s]: failed to form bot: %v", topic, err)
			continue
		}

		// send message
		if payload == nil {
			continue
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
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("webhook error"))
		return
	}
	_, ok := value.String("value")
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("webhook error"))
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
	rw.Write([]byte("ok"))
}
