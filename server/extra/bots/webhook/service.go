package webhook

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/pkg/event"
	"github.com/tinode/chat/server/extra/pkg/route"
	extraStore "github.com/tinode/chat/server/extra/store"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"io"
)

const serviceVersion = "v1"

func webhook(req *restful.Request, resp *restful.Response) {
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

	uid, _ := extraTypes.KV(p.Params).String("uid")
	userUid := types.ParseUserId(uid)
	botUid, _, _, _, err := store.Users.GetAuthUniqueRecord("basic", fmt.Sprintf("%s_bot", Name))
	if err != nil {
		route.ErrorResponse(resp, "error bot")
		return
	}
	topic := userUid.P2PName(botUid)

	d, _ := io.ReadAll(req.Request.Body)

	txt := ""
	if len(d) > 1000 {
		txt = fmt.Sprintf("[webhook:%s] body too long", flag)
	} else {
		txt = fmt.Sprintf("[webhook:%s] %s", flag, string(d))
	}
	// send
	err = event.Emit(event.SendEvent, extraTypes.KV{
		"topic":     topic,
		"topic_uid": int64(botUid),
		"message":   txt,
	})
	if err != nil {
		logs.Err.Println(err)
		_, _ = resp.Write([]byte("send error"))
		return
	}

	_, _ = resp.Write([]byte("ok"))
}
