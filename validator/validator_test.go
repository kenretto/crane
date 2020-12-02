package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"golang.org/x/text/language"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var v = New(Translator{en.New()}, Translator{zh.New()})
	type Test struct {
		RequiredString   string    `json:"required_string" validate:"required"`
		RequiredNumber   int       `json:"required_number" validate:"required"`
		RequiredMultiple []string  `json:"required_multiple" validate:"required"`
		LenString        string    `json:"len_string" validate:"len=1"`
		LteTime          time.Time `json:"lte_time" validate:"lte"`
		GtString         string    `json:"gt_string" validate:"gt=3"`
		GtNumber         float64   `json:"gt_number" validate:"gt=5.56"`
		GtMultiple       []string  `json:"gt_multiple" validate:"gt=2"`
		GtTime           time.Time `json:"gt_time" validate:"gt"`
		GteString        string    `json:"gte_string" validate:"gte=3"`
		GteNumber        float64   `json:"gte_number" validate:"gte=5.56"`
		GteMultiple      []string  `json:"gte_multiple" validate:"gte=2"`
		GteTime          time.Time `json:"gte_time" validate:"gte"`
		EqFieldString    string    `json:"eq_field_string" validate:"eqfield=MaxString"`
		EqCSFieldString  string    `json:"eq_cs_field_string" validate:"eqcsfield=Inner.EqCSFieldString"`
		NeCSFieldString  string    `json:"ne_cs_field_string" validate:"necsfield=Inner.NeCSFieldString"`
	}

	for s, s2 := range v.ValidateStruct(Test{}).(ValidationErrors).Translate(language.Chinese.String()) {
		t.Log(s, s2)
	}
}

// HTTPTestServer http 接口测试服务
type HTTPTestServer struct {
	eng *gin.Engine
}

// NewHTTPTestServer http 接口测试服务
func NewHTTPTestServer() *HTTPTestServer {
	gin.SetMode(gin.ReleaseMode)
	var eng = gin.New()
	binding.Validator = New(Translator{en.New()}, Translator{zh.New()})
	return &HTTPTestServer{
		eng: eng,
	}
}

// HTTPRequest 接口测试
func (httpTestServer *HTTPTestServer) HTTPRequest(req *http.Request, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	httpTestServer.eng.Handle(req.Method, req.URL.Path, handler)
	w := httptest.NewRecorder()
	httpTestServer.eng.ServeHTTP(w, req)
	return w
}

func TestGin(t *testing.T) {
	type Test struct {
		RequiredString   string    `json:"required_string" validate:"required"`
		RequiredNumber   int       `json:"required_number" validate:"required"`
		RequiredMultiple []string  `json:"required_multiple" validate:"required"`
		LenString        string    `json:"len_string" validate:"len=1"`
		LteTime          time.Time `json:"lte_time" validate:"lte"`
		GtString         string    `json:"gt_string" validate:"gt=3"`
		GtNumber         float64   `json:"gt_number" validate:"gt=5.56"`
		GtMultiple       []string  `json:"gt_multiple" validate:"gt=2"`
		GtTime           time.Time `json:"gt_time" validate:"gt"`
		GteString        string    `json:"gte_string" validate:"gte=3"`
		GteNumber        float64   `json:"gte_number" validate:"gte=5.56"`
		GteMultiple      []string  `json:"gte_multiple" validate:"gte=2"`
		GteTime          time.Time `json:"gte_time" validate:"gte"`
		EqFieldString    string    `json:"eq_field_string" validate:"eqfield=MaxString"`
		EqCSFieldString  string    `json:"eq_cs_field_string" validate:"eqcsfield=Inner.EqCSFieldString"`
		NeCSFieldString  string    `json:"ne_cs_field_string" validate:"necsfield=Inner.NeCSFieldString"`
	}
	var req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	response := NewHTTPTestServer().HTTPRequest(req, func(context *gin.Context) {
		var test Test
		var err = binding.Validator.(*Validator).Bind(context, &test)
		if !err.IsValid() {
			context.JSON(http.StatusOK, err.ErrorsInfo)
			return
		}

		context.JSON(http.StatusOK, "false")
	})

	t.Log(response.Body.String())
}