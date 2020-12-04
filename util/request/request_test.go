package request

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestDo(t *testing.T) {
	l := logrus.New()
	logrus.SetReportCaller(true)
	logger = l.WithField("filter", "pkg.util.request.test")
	response, err := URL("http://www.baidu.com/").ContentType("application/json").Parameters(url.Values{"nickname": []string{"wang"}}).Body([]byte(`{"hello": "world"}`)).Do()
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(body))
}
