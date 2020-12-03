package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/sessions"
	"net/http"
)

func SessionSet(ctx *gin.Context) {
	sessions.Set(ctx, "hello", ctx.Query("content"))
	ctx.String(http.StatusOK, "set success")
}

func SessionGet(ctx *gin.Context) {
	ctx.String(http.StatusOK, sessions.Get(ctx, "hello"))
}

func SessionDel(ctx *gin.Context) {
	sessions.Del(ctx, "hello")
	ctx.String(http.StatusOK, "del success")
}
