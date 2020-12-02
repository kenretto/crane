package response

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/i18n"
	"net/http"
)

var (
	Success          = &Response{Code: 2000, Message: "success"}
	PermissionDenied = &Response{Code: 4003, Message: "permission denied"}
	NotFound         = &Response{Code: 4004, Message: "not found"}
	ManyRequest      = &Response{Code: 4429, Message: "too many requests"}
	Failed           = &Response{Code: 5000, Message: "failed"}
	AuthFail         = &Response{Code: 4001, Message: "auth failed"}
)

var (
	// 翻译
	translate *i18n.Bundle
)

// SetTranslator 设置翻译对象
func SetTranslator(t *i18n.Bundle) {
	translate = t
}

// Translator 翻译, 通过分析 AcceptLanguage 来获取用户接受的语言
func Translator(ctx *gin.Context) *i18n.Printer {
	return translate.NewPrinter(i18n.GetAcceptLanguages(ctx)...)
}

// Response HTTP返回数据结构体, 可使用这个, 也可以自定义
type Response struct {
	Code    int         `json:"code"`    // 状态码,这个状态码是与前端和APP约定的状态码,非HTTP状态码
	Data    interface{} `json:"data"`    // 返回数据
	Message string      `json:"message"` // 自定义返回的消息内容
	msgData interface{} // 消息解析使用的数据
}

func (rsp *Response) JSON(data interface{}) *Response {
	rsp.Data = data
	return rsp
}

// Msg msg 描述
func (rsp *Response) Msg(msg string) *Response {
	rsp.Message = msg
	return rsp
}

// MsgData 消息解析使用的数据
func (rsp *Response) MsgData(data interface{}) *Response {
	rsp.msgData = data
	return rsp
}

// End 在调用了这个方法之后,还是需要 return 的
func (rsp *Response) End(c *gin.Context, httpStatus ...int) {
	status := http.StatusOK
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}

	if translate == nil {
		rsp.Message = "please call response.SetTranslator first"
	} else {
		if rsp.Message != "" {
			rsp.Message = translate.NewPrinter(i18n.GetAcceptLanguages(c)...).Translate(rsp.Message, rsp.msgData)
		}
	}

	c.JSON(status, rsp)
}

// Object 直接获得本对象
func (rsp *Response) Object(_ *gin.Context) *Response {
	return rsp
}

// NewResponse 接口返回统一使用这个
//  code 服务端与客户端和web端约定的自定义状态码
//  data 具体的返回数据
//  message 可不传,自定义消息
func NewResponse(code int, data interface{}, message ...string) *Response {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	return &Response{Code: code, Data: data, Message: msg}
}
