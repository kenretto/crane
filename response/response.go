package response

import (
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/i18n"
	"net/http"
)

var (
	// Success success response
	Success = &Response{Code: 2000, Message: "success"}
	// PermissionDenied permission denied response
	PermissionDenied = &Response{Code: 4003, Message: "permission denied"}
	// NotFound not found response
	NotFound = &Response{Code: 4004, Message: "not found"}
	// ManyRequest many request response
	ManyRequest = &Response{Code: 4429, Message: "too many requests"}
	// Failed failed response
	Failed = &Response{Code: 5000, Message: "failed"}
	// AuthFail auth fail response
	AuthFail = &Response{Code: 4001, Message: "auth failed"}
)

var (
	// translator
	translate *i18n.Bundle
)

// SetTranslator set translation object
func SetTranslator(t *i18n.Bundle) {
	translate = t
}

// Translator The language accepted by users can be obtained by analyzing AcceptLanguages
func Translator(ctx *gin.Context) *i18n.Printer {
	return translate.NewPrinter(i18n.GetAcceptLanguages(ctx)...)
}

// Response HTTP return the data structure. You can use this or customize it
type Response struct {
	Code    int         `json:"code"`    // the status code is the status code agreed with the front end and app, not the HTTP status code
	Data    interface{} `json:"data"`    // return data
	Message string      `json:"message"` // customize the returned message content
	msgData interface{} // data used for message parsing
}

// JSON set Data, will return json
func (rsp *Response) JSON(data interface{}) *Response {
	rsp.Data = data
	return rsp
}

// Msg msg description
func (rsp *Response) Msg(msg string) *Response {
	rsp.Message = msg
	return rsp
}

// MsgData data used for message parsing
func (rsp *Response) MsgData(data interface{}) *Response {
	rsp.msgData = data
	return rsp
}

// End after calling this method, you still need return
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

// Object get the object directly
func (rsp *Response) Object(_ *gin.Context) *Response {
	return rsp
}

// NewResponse response
//  code custom status code agreed by server, client and Web
//  data specific return data
//  message can not be transmitted, custom message
func NewResponse(code int, data interface{}, message ...string) *Response {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	return &Response{Code: code, Data: data, Message: msg}
}
