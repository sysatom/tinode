package route

import (
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
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

	// CORS
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		CookiesAllowed: false,
		Container:      restfulContainer}
	restfulContainer.Filter(cors.Filter)

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

func WebService(group, version string, rs ...Router) *restful.WebService {
	path := "/" + prefix + "/" + group + "/" + version
	ws := newWebService(group, version)
	for _, router := range rs {
		funcName := utils.GetFunctionName(router.Function)
		_, resource := utils.ParseFunctionName(funcName)
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
		logs.Info.Printf("WebService %s \t%s%s -> %s", router.Method, path, router.Path, funcName)
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
