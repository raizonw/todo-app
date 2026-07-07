package web_transport_http

import (
	"net/http"
	"os"
	"path"

	core_http_server "github.com/raizonw/todo-app/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHTTPHandler(
	webService WebService,
) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	publicPath := path.Join(os.Getenv("PROJECT_ROOT"), "web", "public")
	fileServer := http.FileServer(http.Dir(publicPath))
	return []core_http_server.Route{
		{
			Path:    "/",
			Handler: fileServer.ServeHTTP,
		},
	}
}
