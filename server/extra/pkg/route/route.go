package route

import (
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/types"
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
	return ws
}

func WebService(group, version string, rs ...*Router) *restful.WebService {
	path := "/" + prefix + "/" + group + "/" + version
	ws := newWebService(group, version)
	for _, router := range rs {
		funcName := utils.GetFunctionName(router.Function)
		_, resource := utils.ParseFunctionName(funcName)
		tags := []string{group}
		var builder *restful.RouteBuilder
		// method
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
		// params
		if len(router.Params) > 0 {
			for _, param := range router.Params {
				switch param.Type {
				case PathParamType:
					builder.Param(ws.PathParameter(param.Name, param.Description).DataType(param.DataType))
				case QueryParamType:
					builder.Param(ws.QueryParameter(param.Name, param.Description).DataType(param.DataType))
				case FormParamType:
					builder.Param(ws.FormParameter(param.Name, param.Description).DataType(param.DataType))
				}
			}
		}
		ws.Route(builder.
			To(router.Function).
			Doc(router.Documentation).
			Operation(router.Method+resource).
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Returns(http.StatusOK, "OK", router.ReturnSample).
			Writes(router.WriteSample))
		logs.Info.Printf("WebService %s \t%s%s \t-> %s", router.Method, path, router.Path, funcName)
	}
	return ws
}

func Route(method string, path string, function restful.RouteFunction, documentation string, options ...Option) *Router {
	r := &Router{
		Method:        method,
		Path:          path,
		Function:      function,
		Documentation: documentation,
		Params:        make([]*Param, 0),
	}
	for _, option := range options {
		option(r)
	}
	return r
}

type Option func(r *Router)

func WithReturns(returns interface{}) Option {
	return func(r *Router) {
		r.ReturnSample = returns
	}
}

func WithWrites(writes interface{}) Option {
	return func(r *Router) {
		r.WriteSample = writes
	}
}

func WithParam(param *Param) Option {
	return func(r *Router) {
		r.Params = append(r.Params, param)
	}
}

type Router struct {
	Method        string
	Path          string
	Function      restful.RouteFunction
	Documentation string
	ReturnSample  interface{}
	WriteSample   interface{}
	Params        []*Param
}

type ParamType string

const (
	PathParamType  ParamType = "path"
	QueryParamType ParamType = "query"
	FormParamType  ParamType = "form"
)

type Param struct {
	Type        ParamType
	Name        string
	Description string
	DataType    string
}

func ErrorResponse(resp *restful.Response, text string) {
	resp.WriteHeader(http.StatusBadRequest)
	_, _ = resp.Write([]byte(text))
}

func URL(group, version string, path string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", types.AppUrl(), prefix, group, version, path)
}
