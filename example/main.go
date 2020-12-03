package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/kenretto/crane"
	"github.com/kenretto/crane/i18n"
	"github.com/kenretto/crane/redis"
	"github.com/kenretto/crane/response"
	"github.com/kenretto/crane/sessions"
	"github.com/kenretto/crane/validator"
	"github.com/kenretto/crudman"
	"github.com/kenretto/crudman/driver"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type member struct {
	ID       int    `gorm:"column:id"`
	Username string `json:"username" form:"username" validate:"min=3,required" gorm:"column:username"`
	Email    string `json:"email" form:"email" validate:"email" gorm:"column:email"`
	Mobile   string `json:"mobile" form:"mobile" validate:"min=6" gorm:"column:mobile"`
	Nickname string `json:"nickname" form:"nickname" validate:"min=5,max=16" gorm:"column:nickname"`
	Avatar   string `json:"avatar" form:"avatar" validate:"url" gorm:"column:avatar"`
}

func (member) TableName() string {
	return "user"
}

type MyCrane struct {
	crane.Crane

	MyRedis *redis.Redis
}

func (my *MyCrane) IntegrationRedis() {
	my.MyRedis = new(redis.Redis)
	my.MyRedis.Config.OnConnect = func(ctx context.Context, cn *redis2.Conn) error {
		fmt.Println("hello, redis")
		return nil
	}
	my.Configurator.Add("redis", my.MyRedis)
}

// Redis get redis
func (my *MyCrane) Redis() *redis.Redis {
	return my.MyRedis
}

func redisMain() {
	var pilot = new(MyCrane)
	err := crane.NewCraneWithCustom("application.yaml", pilot)
	if err != nil {
		panic(err)
	}
	pilot.Handler(func(router *gin.Engine) {
		router.GET("/redis/set", func(context *gin.Context) {
			pilot.Redis().Instance().Set(context, "redis", "yes", 0)
		})
		router.GET("/redis/get", func(context *gin.Context) {
			context.String(http.StatusOK, pilot.Redis().Instance().Get(context, "redis").String())
		})
	})
	pilot.Run()
}

func TestController(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

func main() {
	var pilot, err = crane.NewCrane("application.yaml")
	if err != nil {
		panic(err)
	}

	var validate = validator.New(validator.Translator{Translator: en.New()}, validator.Translator{Translator: zh.New()})
	binding.Validator = validate
	pilot.Handler(func(router *gin.Engine) {
		router.GET("/", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "hello, pilot")
		})

		s := pilot.Sessions().Inject(router)
		s.GET("/session/set", func(ctx *gin.Context) {
			sessions.Set(ctx, "hello", ctx.Query("content"))
			ctx.String(http.StatusOK, "set success")
		})
		s.GET("/session/get", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, sessions.Get(ctx, "hello"))
		})
		s.GET("/session/del", func(ctx *gin.Context) {
			sessions.Del(ctx, "hello")
			ctx.String(http.StatusOK, "del success")
		})
		s.GET("/captcha/img", func(ctx *gin.Context) {
			id, b64s, err := pilot.Captcha().Generate()
			if err != nil {
				pilot.Logger().Error(err)
			}
			sessions.Set(ctx, "captcha_id", id)
			ctx.Header("Content-Type", "image/png")
			data, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(b64s, "data:image/png;base64,", ""))
			if err != nil {
				pilot.Logger().Error(err)
			}
			_, _ = ctx.Writer.Write(data)
		})
		s.GET("captcha/verify", func(ctx *gin.Context) {
			b, _ := strconv.ParseBool(ctx.Query("clear"))
			if pilot.Captcha().Verify(sessions.Get(ctx, "captcha_id"), ctx.Query("answer"), b) {
				ctx.String(http.StatusOK, "ok")
				return
			}
			ctx.String(http.StatusOK, "false")
		})
		s.POST("member/register", func(ctx *gin.Context) {
			var m member
			response.SetTranslator(i18n.NewBundle(language.Chinese).LoadFiles("testdata/i18n/", "yaml", yaml.Unmarshal))
			err := validate.Bind(ctx, &m)
			if !err.IsValid() {
				response.Failed.JSON(err.ErrorsInfo).End(ctx)
				return
			}

			var ormError = pilot.ORM().Save(&m).Error
			if ormError != nil {
				response.Failed.Msg("create error").JSON(ormError).End(ctx)
				return
			}

			response.Success.End(ctx)
		})
		var crud = crudman.New()

		crud.Register(driver.NewGorm(pilot.ORM(), "ID").WithValidator(func(obj interface{}) interface{} {
			return validate.ValidateStruct(obj)
		}), member{}, crudman.SetRoute("/member"))
		response.SetTranslator(i18n.NewBundle(language.Chinese).LoadFiles("testdata/i18n/", "yaml", yaml.Unmarshal))
		s.Any("/crud/*any", func(context *gin.Context) {
			context.Request.URL.Path = strings.ReplaceAll(context.Request.URL.Path, "/crud", "")
			context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(context.MustGet(gin.BodyBytesKey).([]byte)))
			data, err := crud.Handler(context.Writer, context.Request)
			if err != nil {
				fmt.Println(reflect.TypeOf(data))
				if e, ok := data.(validator.ValidationErrors); ok {
					response.Failed.JSON(e.Translate()).End(context)
				} else {
					response.Failed.Msg(err.Error()).End(context)
				}
				return
			}

			response.Success.JSON(data).End(context)
		})
	})

	pilot.Run()
}
