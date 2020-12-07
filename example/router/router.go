package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/example/bootstrap"
	"github.com/kenretto/crane/example/controller"
	"github.com/kenretto/crane/example/controller/member"
	"github.com/kenretto/crane/example/model"
	"github.com/kenretto/crudman"
	"net/http"
)

func Router(router *gin.Engine) {
	router.LoadHTMLGlob("web/**")

	var crud = controller.NewCRUD("/crud")
	crud.Register(bootstrap.Pilot().ORM(), bootstrap.Validator(), model.Member{}, crudman.SetRoute("/member"))
	router.Any("/crud/*any", crud.Controller)

	router.POST("member/register", member.Register)
	//router.GET("/tmp", func(ctx *gin.Context) {
	//	ctx.HTML(http.StatusOK, "index", gin.H{
	//		"links": gin.H{
	//			"/captcha/verify": "captcha-verify",
	//			"/session/set":    "set-session",
	//			"/session/get":    "get-session",
	//			"/session/del":    "delete-session",
	//			"/captcha/img":    "captcha-image",
	//			"/crud/member":    "list-member",
	//		},
	//	})
	//})
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index", gin.H{
			"links": gin.H{
				"/captcha/verify": "captcha-verify",
				"/session/set":    "set-session",
				"/session/get":    "get-session",
				"/session/del":    "delete-session",
				"/captcha/img":    "captcha-image",
				"/crud/member":    "list-member",
			},
		})
	})

	s := bootstrap.Pilot().Sessions().Inject(router)
	s.GET("/session/set", controller.SessionSet)
	s.GET("/session/get", controller.SessionGet)
	s.GET("/session/del", controller.SessionDel)

	s.GET("captcha/verify", controller.CaptchaVerify)
	s.GET("/captcha/img", controller.CaptchaImg)
}
