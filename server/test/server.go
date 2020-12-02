package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/configurator"
	"github.com/kenretto/crane/server"
	"net/http"
	"time"
)

func main() {
	var s = server.NewHTTPServer(server.NewDefaultLogger(server.Fields{"server": "test"}))
	var c, err = configurator.NewConfigurator("testdata/server.yaml")
	if err != nil {
		panic(err)
	}

	c.Add("server", s)
	s.Handler(func(router *gin.Engine) {
		router.GET("/", func(context *gin.Context) {
			time.Sleep(time.Second * 10)
			context.String(http.StatusOK, fmt.Sprintf("hello world: %s", s.Addr))
		})
	})

	s.Run()
}
