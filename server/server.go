package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kenretto/crane/util/stack"
	"github.com/kenretto/daemon"
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
	PID                  string `mapstructure:"pid"`
	ServiceName          string `mapstructure:"name"`
	Addr                 string `mapstructure:"addr"`
	ShutdownWaitDuration string `mapstructure:"shutdown_wait_duration"`
	GinMode              string `mapstructure:"gin_mode"`

	logger   ILogger
	handlers []func(router *gin.Engine)
	handler  *Handler
	server   *http.Server
	rw       sync.RWMutex

	changed chan struct{}
}

// OnChange When the configuration file changes, the service will be listened again
func (httpServer *HTTPServer) OnChange(viper *viper.Viper) {
	httpServer.rw.Lock()
	_ = viper.Unmarshal(httpServer)
	if httpServer.server != nil {
		httpServer.logger.Info("server config changed, re-listening")
	}
	daemon.Register(daemon.NewProcess(httpServer))
	httpServer.rw.Unlock()
	httpServer.changed <- struct{}{}
}

// NewHTTPServer simply initialize the HTTP server
func NewHTTPServer(logger ILogger) *HTTPServer {
	var s = &HTTPServer{
		logger:   logger,
		handlers: make([]func(router *gin.Engine), 0),
		changed:  make(chan struct{}),
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

// PidSavePath get pid save path
func (httpServer *HTTPServer) PidSavePath() string {
	return httpServer.PID
}

// Name get service name
func (httpServer *HTTPServer) Name() string {
	return httpServer.ServiceName
}

// SetLogger set custom logger
func (httpServer *HTTPServer) SetLogger(logger ILogger) {
	httpServer.rw.Lock()
	defer httpServer.rw.Unlock()
	httpServer.logger = logger
}

// Run start server listen, In fact, it is directly called globally daemon.Run() is an effect
func (httpServer *HTTPServer) Run() {
	if rs := daemon.Run(); rs != nil {
		httpServer.logger.Fatalln(rs)
	}
}

func (httpServer *HTTPServer) do() {
	for range httpServer.changed {
		err := httpServer.Stop()
		if err != nil {
			httpServer.logger.Error(err)
		}
		httpServer.rw.Lock()
		var handler = NewHandler(httpServer.GinMode, httpServer.GINRecovery, httpServer.GINLogger)
		handler.Register(httpServer.handlers...)
		httpServer.handler = handler
		httpServer.server = &http.Server{Handler: httpServer.handler.router, Addr: httpServer.Addr}
		httpServer.rw.Unlock()
		go func() {
			err := httpServer.server.ListenAndServe()
			if err == http.ErrServerClosed {
				httpServer.logger.Info(fmt.Sprintf("service [%s] closed at %s", httpServer.ServiceName, time.Now().Format(time.RFC3339)))
			} else {
				httpServer.logger.Error(fmt.Sprintf("service [%s] error: %v", httpServer.ServiceName, err))
			}
			httpServer.logger.Info(fmt.Sprintf("server starting, listen: %s", httpServer.Addr))
			httpServer.handler.Print(httpServer.logger)
		}()
	}
}

// Start daemon start handle
func (httpServer *HTTPServer) Start() {
	httpServer.logger.Info(fmt.Sprintf("server starting, listen: %s", httpServer.Addr))
	httpServer.handler.Print(httpServer.logger)
}

// Stop daemon  stop handler
func (httpServer *HTTPServer) Stop() error {
	httpServer.rw.Lock()
	defer httpServer.rw.Unlock()
	if httpServer.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpServer.shutdownDuration())
	defer cancel()
	err := httpServer.server.Shutdown(ctx)
	httpServer.logger.Info(fmt.Sprintf("service [%s] stoped", httpServer.ServiceName))
	return err
}

// Restart daemon restart
func (httpServer *HTTPServer) Restart() error {
	httpServer.logger.Info(fmt.Sprintf("service [%s] restarting", httpServer.ServiceName))
	return httpServer.Stop()
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
	params["latency"] = time.Since(start)
	params["method"] = ctx.Request.Method
	params["status"] = ctx.Writer.Status()
	params["body_size"] = ctx.Writer.Size()
	params["body"] = request
	params["client_ip"] = ctx.ClientIP()
	params["user_agent"] = ctx.Request.UserAgent()
	params["log_type"] = "pkg.server.server"
	params["keys"] = ctx.Keys
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
			httpRequest, _ := httputil.DumpRequest(ctx.Request, false)

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
