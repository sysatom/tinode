package route

import (
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/utils"
	"net/http"
)

const prefix = "bot"

func NewContainer() *restful.Container {
	restfulContainer := restful.NewContainer()
	restfulContainer.ServeMux = http.NewServeMux()
	restfulContainer.Router(restful.CurlyRouter{})

	restfulContainer.RecoverHandler(func(panicReason interface{}, w http.ResponseWriter) {
		logStackOnRecover(panicReason, w)
	})
	restfulContainer.ServiceErrorHandler(func(serviceError restful.ServiceError, req *restful.Request, resp *restful.Response) {
		logServiceError(serviceError, req, resp)
	})

	return restfulContainer
}

func newWebService(group string, version string) *restful.WebService {
	ws := new(restful.WebService)
	path := "/" + prefix + "/" + group + "/" + version
	ws.Path(path)
	ws.Doc(fmt.Sprintf("API at %s", path))
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	return ws
}

func RegisterResourceHandlers(ws *restful.WebService, group string, function restful.RouteFunction) { // fixme group
	resource := group
	tags := []string{resource}
	ws.Route(ws.GET("/user/{id}").To(function).
		Doc(fmt.Sprintf("Get %s resource", resource)).
		Operation(resource+"Get").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "OK", model.Objective{}).
		Writes(model.Objective{}))
}

func WebService(group, version string, rs ...Router) *restful.WebService {
	ws := newWebService(group, version)
	for _, router := range rs {
		resource := utils.GetFunctionName(router.Function)
		tags := []string{resource}
		var builder *restful.RouteBuilder
		switch router.Method {
		case "GET":
			builder = ws.GET(router.Path)
		case "POST":
			builder = ws.POST(router.Path)
		case "PUT":
			builder = ws.PUT(router.Path)
		case "PATCH":
			builder = ws.PATCH(router.Path)
		case "DELETE":
			builder = ws.DELETE(router.Path)
		default:
			continue
		}
		ws.Route(builder.
			To(router.Function).
			Doc(router.Documentation).
			Operation(router.Method+resource).
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Returns(http.StatusOK, "OK", router.ReturnSample).
			Writes(router.WriteSample))
	}
	return ws
}

func Route(method string, path string, function restful.RouteFunction, documentation string, returns, writes interface{}) Router {
	return Router{
		Method:        method,
		Path:          path,
		Function:      function,
		Documentation: documentation,
		ReturnSample:  returns,
		WriteSample:   writes,
	}
}

type Router struct {
	Method        string
	Path          string
	Function      restful.RouteFunction
	Documentation string
	ReturnSample  interface{}
	WriteSample   interface{}
}
