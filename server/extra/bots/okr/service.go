package okr

import (
	"embed"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/store/types"
	"net/http"
	"strconv"
)

const serviceVersion = "v1"

//go:embed webapp/build
var dist embed.FS

func webapp(rw http.ResponseWriter, req *http.Request) {
	bots.ServeFile(rw, req, dist, "webapp/build")
}

func objectiveList(req *restful.Request, resp *restful.Response) {
	uid, _ := req.Attribute("uid").(types.Uid)
	topic, _ := req.Attribute("uid").(string)
	list, err := store.Chatbot.ListObjectives(uid, topic)
	if err != nil {
		_ = resp.WriteAsJson(extraTypes.ErrMessage(http.StatusBadRequest, err.Error()))
		return
	}
	_ = resp.WriteAsJson(extraTypes.OkMessage(list))
}

func objectiveDetail(req *restful.Request, resp *restful.Response) {
	uid, _ := req.Attribute("uid").(types.Uid)
	topic, _ := req.Attribute("uid").(string)
	s := req.PathParameter("sequence")
	sequence, _ := strconv.ParseInt(s, 10, 64)

	obj, err := store.Chatbot.GetObjectiveBySequence(uid, topic, sequence)
	if err != nil {
		_ = resp.WriteAsJson(extraTypes.ErrMessage(http.StatusNotFound, ""))
		return
	}
	_ = resp.WriteAsJson(extraTypes.OkMessage(obj))
}

func objectiveCreate(req *restful.Request, resp *restful.Response) {
	uid, _ := req.Attribute("uid").(types.Uid)
	topic, _ := req.Attribute("uid").(string)
	obj := new(model.Objective)
	err := req.ReadEntity(&obj)
	if err != nil {
		_ = resp.WriteAsJson(extraTypes.ErrMessage(http.StatusNotFound, err.Error()))
		return
	}
	obj.UID = uid.UserId()
	obj.Topic = topic
	_, err = store.Chatbot.CreateObjective(obj)
	if err != nil {
		_ = resp.WriteAsJson(extraTypes.ErrMessage(http.StatusNotFound, err.Error()))
		return
	}
	_ = resp.WriteAsJson(extraTypes.OkMessage(nil))
}

func objectiveUpdate(req *restful.Request, resp *restful.Response) {
	uid, _ := req.Attribute("uid").(types.Uid)
	topic, _ := req.Attribute("uid").(string)
	obj := new(model.Objective)
	err := req.ReadEntity(&obj)
	if err != nil {
		_ = resp.WriteAsJson(extraTypes.ErrMessage(http.StatusNotFound, err.Error()))
		return
	}
	obj.UID = uid.UserId()
	obj.Topic = topic
	err = store.Chatbot.UpdateObjective(obj)
	if err != nil {
		_ = resp.WriteAsJson(extraTypes.ErrMessage(http.StatusNotFound, err.Error()))
		return
	}
	_ = resp.WriteAsJson(extraTypes.OkMessage(nil))
}

func objectiveDelete(req *restful.Request, resp *restful.Response) {
	uid, _ := req.Attribute("uid").(types.Uid)
	topic, _ := req.Attribute("uid").(string)
	s := req.PathParameter("sequence")
	sequence, _ := strconv.ParseInt(s, 10, 64)

	err := store.Chatbot.DeleteObjectiveBySequence(uid, topic, sequence)
	if err != nil {
		_ = resp.WriteAsJson(extraTypes.ErrMessage(http.StatusNotFound, err.Error()))
		return
	}
	_ = resp.WriteAsJson(extraTypes.OkMessage(nil))
}
