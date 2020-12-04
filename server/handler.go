package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strings"
)

type debugPrintRouteInfo struct {
	handler, route string
	numHandlers    int
}

// Handler router
type Handler struct {
	router              *gin.Engine
	debugPrintRouteInfo []debugPrintRouteInfo
}

// Register register handler
func (handler *Handler) Register(handlers ...func(router *gin.Engine)) {
	for _, f := range handlers {
		f(handler.router)
	}
}

// Group router group
func (*Handler) Group(relativePath string, handlers ...func(router gin.IRouter)) func(router gin.IRouter) {
	return func(router gin.IRouter) {
		var group = router.Group(relativePath)
		for _, f := range handlers {
			f(group)
		}
	}
}

// NewHandler init handler
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

// Print print route info
func (handler *Handler) Print() {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:  true,
		DisableQuote: true,
	})

	var handlerMaxLength = 0
	for _, info := range handler.router.Routes() {
		if len(info.Handler) > handlerMaxLength {
			handlerMaxLength = len(info.Handler)
		}
	}

	for _, info := range handler.router.Routes() {
		var file, line = runtime.FuncForPC(reflect.ValueOf(info.HandlerFunc).Pointer()).FileLine(reflect.ValueOf(info.HandlerFunc).Pointer())
		logger.WithField("handler", info.Handler+strings.Repeat(" ", handlerMaxLength+4-len(info.Handler))+fmt.Sprintf("%s:%d", file, line)).
			Println(info.Method + strings.Repeat(" ", 10-len(info.Method)) + info.Path)
	}
}

// Router get *gin.Engine
func (handler *Handler) Router() *gin.Engine {
	return handler.router
}
