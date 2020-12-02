package middleware

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/response"
	"github.com/kenretto/crane/sessions"
	"math/rand"
	"net/http"
	"time"
)

type CSRFToken struct {
	Domain     string `mapstructure:"domain"`
	Length     int    `mapstructure:"length"`
	SessionKey string `mapstructure:"session_key"`
	Duration   string `mapstructure:"duration"`
	Path       string `mapstructure:"path"`
	Secure     bool   `mapstructure:"secure"`
	HTTPOnly   bool   `mapstructure:"http_only"`
}

func (c CSRFToken) duration() time.Duration {
	d, err := time.ParseDuration(c.Duration)
	if err != nil {
		return time.Hour * 2
	}
	return d
}

func (c CSRFToken) Token() string {
	rand.Seed(time.Now().UnixNano())
	var token = make([]byte, c.Length)
	for i := 0; i < c.Length; i++ {
		token[i] = byte(rand.Intn(127))
	}

	return base64.URLEncoding.EncodeToString(token)
}

func (c CSRFToken) Valid(ctx *gin.Context) error {
	token, err := ctx.Cookie(c.SessionKey)
	if err != nil || token == "" {
		return errors.New("csrf token error")
	}
	if token != sessions.Get(ctx, c.SessionKey) {
		return errors.New("csrf token error")
	}

	return nil
}

func (c CSRFToken) Handler(ctx *gin.Context) {
	switch ctx.Request.Method {
	case http.MethodGet:
		var token = c.Token()
		sessions.Del(ctx, c.SessionKey)
		sessions.Set(ctx, c.SessionKey, token)
		ctx.SetCookie(c.SessionKey, token, int(c.duration()/time.Second), c.Path, c.Domain, c.Secure, c.HTTPOnly)
	default:
		if err := c.Valid(ctx); err != nil {
			response.Failed.Msg(err.Error()).End(ctx)
			ctx.Abort()
		}
		sessions.Del(ctx, c.SessionKey)
		ctx.Next()
	}
}
