package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/kenretto/crane"
	"github.com/kenretto/crane/example/bootstrap"
	"github.com/kenretto/crane/example/router"
	"github.com/kenretto/crane/redis"
	"net/http"
)

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
	bootstrap.Pilot().Handler(router.Router)
	bootstrap.Pilot().Run()
}
