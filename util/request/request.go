// Package request 本项目对外发起请求使用此包
package request

import (
	"bytes"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	client = http.DefaultClient
	logger *logrus.Entry
)

type (
	// Request 请求数据
	Request struct {
		method      string
		contentType string
		headers     http.Header
		url         string
		parameters  url.Values
		body        []byte
	}

	// Response 扩展 *http.Response 的功能
	Response struct {
		*http.Response
	}
)

// SetLogger 设置全局 log 输出对象
func SetLogger(l *logrus.Entry) {
	logger = l
}

// Unmarshal 将 response 数据当作 json 转换为结构体或map
func (response *Response) Unmarshal(v interface{}) error {
	var body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(body, v)
}

// URL 构造请求
// 默认使用 GET 请求
func URL(url string) *Request {
	return &Request{url: url, method: http.MethodGet}
}

// Method 指定请求方式
func (req *Request) Method(method string) *Request {
	req.method = method
	return req
}

// ContentType 指定传参类型
func (req *Request) ContentType(contentType string) *Request {
	req.contentType = contentType
	return req
}

// Header 指定header
func (req *Request) Header(header http.Header) *Request {
	req.headers = header
	return req
}

// Parameters 指定 url 参数
func (req *Request) Parameters(parameters url.Values) *Request {
	req.parameters = parameters
	return req
}

// Body 指定 body 参数
func (req *Request) Body(body []byte) *Request {
	req.body = body
	return req
}

// HTTPRequest 构造 *http.Request 对象
func (req *Request) HTTPRequest() (*http.Request, error) {
	switch req.method {
	case http.MethodHead, http.MethodDelete, http.MethodPatch, http.MethodGet:
		var httpRequest, err = http.NewRequest(req.method, req.url, bytes.NewReader(req.body))
		if err != nil {
			return nil, err
		}
		httpRequest.URL.RawQuery = req.parameters.Encode()
		httpRequest.Header = req.headers
		return httpRequest, nil
	case http.MethodPost, http.MethodPut:
		var (
			httpRequest *http.Request
			err         error
		)
		switch req.contentType {
		case "application/json":
			httpRequest, err = http.NewRequest(req.method, req.url, bytes.NewReader(req.body))
			if err != nil {
				return nil, err
			}
		default:
			httpRequest, err = http.NewRequest(req.method, req.url, strings.NewReader(req.parameters.Encode()))
			if err != nil {
				return nil, err
			}
		}
		httpRequest.Header.Set("Content-Type", req.contentType)
		return httpRequest, nil
	}
	return nil, errors.New("unsupported method")
}

// Do 发起请求
func (req *Request) Do() (*Response, error) {
	return Do(req)
}

// Do 发起请求
func Do(req *Request) (*Response, error) {
	var (
		fields = logrus.Fields{
			"parameters":     req.parameters,
			"request_method": req.method,
			"request_url":    req.url,
			"request_body":   string(req.body),
			"content-type":   req.contentType,
		}
		begin = time.Now()
	)

	request, err := req.HTTPRequest()
	if err != nil {
		return nil, err
	}
	client.Timeout = time.Second * 30
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	fields["duration"] = time.Since(begin).String()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	fields["response_body"] = string(body)
	logger.WithFields(fields).Info()

	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return &Response{response}, err
}
