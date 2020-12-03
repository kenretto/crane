package controller

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/example/bootstrap"
	"github.com/kenretto/crane/sessions"
	"net/http"
	"strconv"
	"strings"
)

func CaptchaVerify(ctx *gin.Context) {
	b, _ := strconv.ParseBool(ctx.Query("clear"))
	if ctx.Query("answer") == "" {
		ctx.String(http.StatusOK, "false")
		return
	}
	if bootstrap.Pilot().Captcha().Verify(sessions.Get(ctx, "captcha_id"), ctx.Query("answer"), b) {
		ctx.String(http.StatusOK, "ok")
		return
	}
	ctx.String(http.StatusOK, "false")
}

func CaptchaImg(ctx *gin.Context) {
	id, b64s, err := bootstrap.Pilot().Captcha().Generate()
	if err != nil {
		bootstrap.Pilot().Logger().Error(err)
	}
	sessions.Set(ctx, "captcha_id", id)
	ctx.Header("Content-Type", "image/png")
	data, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(b64s, "data:image/png;base64,", ""))
	if err != nil {
		bootstrap.Pilot().Logger().Error(err)
	}
	_, _ = ctx.Writer.Write(data)
}
