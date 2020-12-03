// Package test test helper
package test

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane"
	"github.com/kenretto/crane/server"
	"net/http"
	"net/http/httptest"
)

// Test common test
type Test struct {
	service crane.ICrane
	handler *server.Handler
}

func NewTest(config string) *Test {
	var (
		test = new(Test)
		err  error
	)
	test.service, err = crane.NewCrane(config)
	if err != nil {
		panic(err)
	}

	test.handler = server.NewHandler(test.service.Server().GinMode, test.service.Server().GINRecovery, test.service.Server().GINLogger)
	return test
}

// HTTPRequest 接口测试
func (test *Test) HTTPRequest(req *http.Request, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	test.handler.Register(func(router *gin.Engine) {
		router.Handle(req.Method, req.URL.Path, handler)
	})
	w := httptest.NewRecorder()
	test.handler.Router().ServeHTTP(w, req)
	return w
}
