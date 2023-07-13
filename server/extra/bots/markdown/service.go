package markdown

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/event"
	"github.com/tinode/chat/server/extra/pkg/route"
	extraStore "github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"text/template"
)

const serviceVersion = "v1"

//go:embed markdown.html
var editorTemplate string

func editor(req *restful.Request, resp *restful.Response) {
	flag := req.PathParameter("flag")

	p, err := extraStore.Chatbot.ParameterGet(flag)
	if err != nil {
		route.ErrorResponse(resp, "flag error")
		return
	}
	if p.IsExpired() {
		route.ErrorResponse(resp, "page expired")
		return
	}

	t, err := template.New("tmpl").Parse(editorTemplate)
	if err != nil {
		route.ErrorResponse(resp, "page template error")
		return
	}
	buf := bytes.NewBufferString("")
	p.Params["flag"] = flag
	data := p.Params
	err = t.Execute(buf, data)

	_, _ = resp.Write(buf.Bytes())
}

func saveMarkdown(req *restful.Request, resp *restful.Response) {
	// data
	var data map[string]string
	err := req.ReadEntity(&data)
	if err != nil {
		route.ErrorResponse(resp, "params error")
		return
	}

	uid, _ := data["uid"]
	flag, _ := data["flag"]
	markdown, _ := data["markdown"]
	if uid == "" || flag == "" || markdown == "" {
		route.ErrorResponse(resp, "params error")
		return
	}

	p, err := extraStore.Chatbot.ParameterGet(flag)
	if err != nil {
		route.ErrorResponse(resp, "flag error")
		return
	}
	if p.IsExpired() {
		route.ErrorResponse(resp, "page expired")
		return
	}

	// store
	userUid := types.ParseUserId(uid)
	botUid, _, _, _, err := store.Users.GetAuthUniqueRecord("basic", fmt.Sprintf("%s_bot", Name))
	topic := userUid.P2PName(botUid)
	payload := bots.StorePage(
		extraTypes.Context{AsUser: userUid, Original: botUid.UserId()},
		model.PageMarkdown, "",
		extraTypes.MarkdownMsg{Raw: markdown})
	message := ""
	if link, ok := payload.(extraTypes.LinkMsg); ok {
		message = link.Url
	}

	// send
	err = event.Emit(event.SendEvent, map[string]interface{}{
		"topic":     topic,
		"topic_uid": int64(botUid),
		"message":   message,
	})
	if err != nil {
		logs.Err.Println(err)
		_, _ = resp.Write([]byte("send error"))
		return
	}

	_, _ = resp.Write([]byte("ok"))
}
