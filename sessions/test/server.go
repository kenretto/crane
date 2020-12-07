package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/configurator"
	"github.com/kenretto/crane/server"
	"github.com/kenretto/crane/sessions"
	"net/http"
)

func main() {
	var s = server.NewHTTPServer(server.NewDefaultLogger(server.Fields{"server": "test"}))
	var c, err = configurator.NewConfigurator("testdata/server.yaml")
	if err != nil {
		panic(err)
	}

	var sess = new(sessions.Sessions)
	c.Add(sess)

	s.Handler(func(router *gin.Engine) {
		r := sess.Inject(router)
		r.GET("/set", func(context *gin.Context) {
			sessions.Set(context, "hello", "world")
		})
		r.GET("/get", func(context *gin.Context) {
			context.String(http.StatusOK, sessions.Get(context, "hello"))
		})
		r.GET("/del", func(context *gin.Context) {
			sessions.Del(context, "hello")
		})
	})

	c.Add(s)
	s.Listen()
}
