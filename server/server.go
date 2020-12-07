package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kenretto/crane/util/stack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"
)

// HTTPServer default http server
type HTTPServer struct {
	Addr                 string `mapstructure:"addr"`
	ShutdownWaitDuration string `mapstructure:"shutdown_wait_duration"`
	GinMode              string `mapstructure:"gin_mode"`
	Metrics              string `mapstructure:"metrics"`

	logger   ILogger
	handlers []func(router *gin.Engine)
	handler  *Handler
	server   *http.Server
	rw       sync.RWMutex

	running, changed chan struct{}
	exit             chan struct{}
}

func (httpServer *HTTPServer) Node() string {
	return "server"
}

// OnChange When the configuration file changes, the service will be listened again
func (httpServer *HTTPServer) OnChange(viper *viper.Viper) {
	httpServer.rw.Lock()
	_ = viper.Unmarshal(httpServer)
	if httpServer.server != nil {
		httpServer.logger.Info("server config changed, re-listening")
	}
	httpServer.rw.Unlock()
	httpServer.changed <- struct{}{}
}

// NewHTTPServer simply initialize the HTTP server
func NewHTTPServer(logger ILogger) *HTTPServer {
	var s = &HTTPServer{
		logger:   logger,
		handlers: make([]func(router *gin.Engine), 0),
		changed:  make(chan struct{}),
		running:  make(chan struct{}),
		exit:     make(chan struct{}),
	}

	go s.do()
	return s
}

// Handler use this method to register the handler of gin
func (httpServer *HTTPServer) Handler(handler func(router *gin.Engine)) {
	httpServer.rw.Lock()
	defer httpServer.rw.Unlock()
	httpServer.handlers = append(httpServer.handlers, handler)
}

func (httpServer *HTTPServer) Router() *gin.Engine {
	return httpServer.handler.router
}

func (httpServer *HTTPServer) shutdownDuration() time.Duration {
	duration, err := time.ParseDuration(httpServer.ShutdownWaitDuration)
	if err != nil {
		return time.Second * 30
	}
	return duration
}

// SetLogger set custom logger
func (httpServer *HTTPServer) SetLogger(logger ILogger) {
	httpServer.rw.Lock()
	defer httpServer.rw.Unlock()
	httpServer.logger = logger
}

func (httpServer *HTTPServer) close() {
	if httpServer.server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpServer.shutdownDuration())
	_ = httpServer.server.Shutdown(ctx)
	httpServer.logger.Info("service stop")
	cancel()
}

func (httpServer *HTTPServer) Stop() {
	httpServer.exit <- struct{}{}
}

func (httpServer *HTTPServer) Listen() {
	go func() {
		for range httpServer.running {
			httpServer.logger.Info(fmt.Sprintf("server starting, listen: %s", httpServer.Addr))
			httpServer.handler.Print()
			err := httpServer.server.ListenAndServe()
			if err == http.ErrServerClosed {
				httpServer.logger.Info(fmt.Sprintf("service closed at %s", time.Now().Format(time.RFC3339)))
			} else {
				httpServer.logger.Error(fmt.Sprintf("service error: %v", err))
			}
		}
	}()

	<-httpServer.exit
	httpServer.close()
}

func (httpServer *HTTPServer) do() {
	for range httpServer.changed {
		httpServer.rw.Lock()
		httpServer.close()
		var handler = NewHandler(httpServer.GinMode, httpServer.GINRecovery, httpServer.GINLogger)
		if httpServer.Metrics != "" {
			handler.Register(func(router *gin.Engine) {
				router.GET(httpServer.Metrics, func(context *gin.Context) {
					promhttp.Handler().ServeHTTP(context.Writer, context.Request)
				})
			})
		}
		handler.Register(httpServer.handlers...)
		httpServer.handler = handler
		httpServer.server = &http.Server{Handler: httpServer.handler.router, Addr: httpServer.Addr}
		httpServer.rw.Unlock()
		httpServer.running <- struct{}{}
	}
}

// GINLogger 自定义的GIN日志处理中间件
func (httpServer *HTTPServer) GINLogger(ctx *gin.Context) {
	start := time.Now()
	Metrics.HTTPRequestURICounter(ctx.Request.URL.Path)
	var request map[string]interface{}
	switch ctx.Request.Method {
	case http.MethodPost, http.MethodPut:
		switch ctx.ContentType() {
		case "application/json":
			_ = ctx.ShouldBindBodyWith(&request, binding.JSON)
		case "application/xml", "text/xml":
			_ = ctx.ShouldBindBodyWith(&request, binding.XML)
		default:
			request = make(map[string]interface{})
			err := ctx.Request.ParseForm()
			if err == nil {
				for k, v := range ctx.Request.Form {
					request[k] = v
				}
			}

			form, err := ctx.MultipartForm()
			if err == nil {
				for k, v := range form.Value {
					request[k] = v
				}
			}
		}
	}

	ctx.Next()

	var params = make(Fields)
	params["latency"] = time.Since(start).String()
	params["method"] = ctx.Request.Method
	params["status"] = ctx.Writer.Status()
	params["body_size"] = ctx.Writer.Size()
	params["body"] = request
	params["client_ip"] = ctx.ClientIP()
	params["user_agent"] = ctx.Request.UserAgent()
	params["keys"] = ctx.Keys
	params["headers"] = ctx.Request.Header
	Metrics.HTTPResponseStatusCounter(ctx.Request.URL.Path, ctx.Writer.Status())
	httpServer.logger.WithFields(params).Info(ctx.Request.URL.String())
}

// GINRecovery gin recovery handler
func (httpServer *HTTPServer) GINRecovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			defer Metrics.RequestPanicCounter(ctx.Request.URL.Path)
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			printStack := stack.Stack(3)
			httpRequest, _ := httputil.DumpRequest(ctx.Request, true)

			if gin.Mode() != gin.ReleaseMode {
				httpServer.logger.Error(string(httpRequest))
				var errors = make([]logrus.Fields, 0)
				for i := 0; i < len(printStack); i++ {
					errors = append(errors, logrus.Fields{
						"func":   printStack[i]["func"],
						"source": printStack[i]["source"],
						"file":   fmt.Sprintf("%s:%d", printStack[i]["file"], printStack[i]["line"]),
					})
				}
				httpServer.logger.WithFields(Fields{"stack": errors}).Error(err)
				if gin.Mode() == gin.DebugMode {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"stack": errors, "message": err})
				}
			} else {
				httpServer.logger.WithFields(Fields{"stack": printStack, "request": string(httpRequest)}).Error(err)
			}

			if brokenPipe {
				_ = ctx.Error(err.(error))
				ctx.Abort()
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	ctx.Next()
}
