package member

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/example/bootstrap"
	"github.com/kenretto/crane/example/model"
	"github.com/kenretto/crane/response"
)

func Register(ctx *gin.Context) {
	var m model.Member
	err := bootstrap.Validator().Bind(ctx, &m)
	if !err.IsValid() {
		response.Failed.JSON(err.ErrorsInfo).End(ctx)
		return
	}

	var ormError = bootstrap.Pilot().ORM().Save(&m).Error
	if ormError != nil {
		response.Failed.Msg("create error").JSON(ormError).End(ctx)
		return
	}

	response.Success.End(ctx)
}
