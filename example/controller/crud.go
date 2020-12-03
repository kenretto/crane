package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/response"
	"github.com/kenretto/crane/validator"
	"github.com/kenretto/crudman"
	"github.com/kenretto/crudman/driver"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

type CRUD struct {
	prefix string
	crud   *crudman.Managers
}

func NewCRUD(prefix string) *CRUD {
	return &CRUD{
		prefix: prefix,
		crud:   crudman.New(),
	}
}

func (crud *CRUD) parseRoute(ctx *gin.Context) {
	ctx.Request.URL.Path = strings.ReplaceAll(ctx.Request.URL.Path, crud.prefix, "")
}

func (crud *CRUD) refillBody(ctx *gin.Context) {
	v, exist := ctx.Get(gin.BodyBytesKey)
	if exist {
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(v.([]byte)))
	}
}

func (crud *CRUD) Register(db *gorm.DB, validate *validator.Validator, entity crudman.Tabler, setups ...crudman.Setup) {
	crud.crud.Register(driver.NewGorm(db, "ID").WithValidator(func(obj interface{}) interface{} {
		return validate.ValidateStruct(obj)
	}), entity, setups...)
}

func (crud *CRUD) Controller(ctx *gin.Context) {
	crud.parseRoute(ctx)
	crud.refillBody(ctx)
	data, err := crud.crud.Handler(ctx.Writer, ctx.Request)
	if err != nil {
		if e, ok := data.(validator.ValidationErrors); ok {
			response.Failed.JSON(e.Translate()).End(ctx)
		} else {
			response.Failed.Msg(err.Error()).End(ctx)
		}
		return
	}

	response.Success.JSON(data).End(ctx)
}
