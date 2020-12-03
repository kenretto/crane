package uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUploader_Uploader(t *testing.T) {
	var eng = gin.New()
	eng.Handle(http.MethodPost, "/test/upload", func(context *gin.Context) {
		var uploader = Uploader{
			FormKey:      "files",
			SaveHandler:  new(DefaultSaveHandler).SetDst("testdata/").SetPrefix("avatar_"),
			AllowedTypes: []string{"jpg"},
			NameFn: func(index int, file *multipart.FileHeader) string {
				return fmt.Sprintf("%d.jpg", index)
			},
			Ctx: context,
		}
		files, err := uploader.SaveAll()
		context.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": files,
			"msg": func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
		})
	})
	w := httptest.NewRecorder()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	fw, err := writer.CreateFormFile("files", "testdata/0B627A7A0A8F14F29D4E33246B637A3C.jpg")
	if err != nil {
		t.Error(err)
	}
	f, _ := os.Open("testdata/0B627A7A0A8F14F29D4E33246B637A3C.jpg")
	_, _ = io.Copy(fw, f)

	fw, err = writer.CreateFormFile("files", "testdata/0B627A7A0A8F14F29D4E33246B637A3C.jpg")
	if err != nil {
		t.Error(err)
	}

	f, _ = os.Open("testdata/0B627A7A0A8F14F29D4E33246B637A3C.jpg")
	_, _ = io.Copy(fw, f)

	_ = writer.Close()
	var request = httptest.NewRequest(http.MethodPost, "/test/upload", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	eng.ServeHTTP(w, request)

	var rs map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &rs)
	if err != nil {
		t.Error(err)
	}

	files, ok := rs["data"].([]interface{})
	if !ok {
		t.Error("response data error")
	}

	if len(files) != 2 {
		t.Error("upload error")
	}

	if files[0].(string) != "testdata/avatar_0.jpg" {
		t.Error("file name error")
	}
}
