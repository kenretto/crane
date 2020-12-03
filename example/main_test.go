package main

import (
	"github.com/kenretto/crane/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTestController(t *testing.T) {
	var req = httptest.NewRequest(http.MethodGet, "/test", nil)
	var record = test.NewTest("application.yaml").HTTPRequest(req, TestController)
	if record.Body.String() != "ok" {
		t.Error("test error")
	}
	t.Log(record.Body.String())
}
