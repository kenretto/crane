package server

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type debugPrintRouteInfo struct {
	handler, route string
	numHandlers    int
}

// Handler 路由
type Handler struct {
	router              *gin.Engine
	debugPrintRouteInfo []debugPrintRouteInfo
}

// Register 注册路由
func (handler *Handler) Register(handlers ...func(router *gin.Engine)) {
	for _, f := range handlers {
		f(handler.router)
	}
}

// Group 路由分组
func (*Handler) Group(relativePath string, handlers ...func(router gin.IRouter)) func(router gin.IRouter) {
	return func(router gin.IRouter) {
		var group = router.Group(relativePath)
		for _, f := range handlers {
			f(group)
		}
	}
}

// NewHandler 初始化 handler
func NewHandler(mode string, recovery, ginLogger gin.HandlerFunc) *Handler {
	gin.SetMode(mode)
	var handler = &Handler{
		router:              gin.New(),
		debugPrintRouteInfo: make([]debugPrintRouteInfo, 0),
	}

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		handler.debugPrintRouteInfo = append(handler.debugPrintRouteInfo, debugPrintRouteInfo{
			handler:     handlerName,
			route:       httpMethod + strings.Repeat(" ", 12-len(httpMethod)) + absolutePath,
			numHandlers: nuHandlers,
		})
	}

	handler.router.Use(recovery, ginLogger)
	return handler
}

func (handler *Handler) Print(l ILogger) {
	for _, info := range handler.router.Routes() {
		l.Println(info.Method + strings.Repeat(" ", 12-len(info.Method)) + info.Path)
	}
}

func (handler *Handler) Router() *gin.Engine {
	return handler.router
}
