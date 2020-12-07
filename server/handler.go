package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

var methodColor = map[string]string{
	http.MethodGet:     "\033[0;32mGET\033[0m",
	http.MethodHead:    "\033[0;34mHEAD\033[0m",
	http.MethodPost:    "\033[0;33mPOST\033[0m",
	http.MethodPut:     "\033[0;36mPUT\033[0m",
	http.MethodPatch:   "\033[0;32mPATCH\033[0m",
	http.MethodDelete:  "\033[0;31mDELETE\033[0m",
	http.MethodConnect: "\033[0;37mCONNECT\033[0m",
	http.MethodOptions: "\033[0;37mOPTIONS\033[0m",
	http.MethodTrace:   "\033[0;37mTRACE\033[0m",
}

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
	var (
		logger           = log.New(os.Stdout, "\033[1;32m[Router]\033[0m"+strings.Repeat(" ", 4), log.Lmsgprefix)
		handlerMaxLength = 0
		routerMaxLength  = 0
		routers          = handler.router.Routes()
	)
	sort.Slice(routers, func(i, j int) bool {
		return routers[i].Path < routers[j].Path
	})

	for _, info := range routers {
		if len(info.Handler) > handlerMaxLength {
			handlerMaxLength = len(info.Handler)
		}

		var paddingText = strings.Repeat(" ", 10-len(info.Method))
		if len(info.Method+paddingText+info.Path) > routerMaxLength {
			routerMaxLength = len(info.Method + paddingText + info.Path)
		}
	}

	for _, info := range routers {
		var (
			handlerPointer = reflect.ValueOf(info.HandlerFunc).Pointer()
			file, line     = runtime.FuncForPC(handlerPointer).FileLine(handlerPointer)
			paddingText    = strings.Repeat(" ", 10-len(info.Method))
			text           = methodColor[info.Method] + paddingText + info.Path + strings.Repeat(" ", routerMaxLength+4-len(info.Method+paddingText+info.Path)) + fmt.Sprintf("\033[0;36m%s\033[0m", info.Handler) + strings.Repeat(" ", handlerMaxLength+4-len(info.Handler)) + fmt.Sprintf("%s:%d", file, line)
		)
		logger.Println(fmt.Sprintf("\033[1;32m%s\033[0m", text))
	}
}

// Router get *gin.Engine
func (handler *Handler) Router() *gin.Engine {
	return handler.router
}
