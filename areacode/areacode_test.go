package areacode

import (
	"testing"
	"time"
)

func TestGetLocalTimeByCode(t *testing.T) {
	var localTime, err = GetLocalTimeByCode(Filter{Code: 81}, "2020-08-18 14:00:00")
	if err != nil {
		t.Fatal(err)
	}

	if localTime-time.Date(2020, time.August, 18, 14, 0, 00, 0, time.Local).Unix() != -3600 {
		t.Errorf("time error: localTime=%d", localTime)
	}

	var t1 = "2020-08-18 12:44:00"
	localTime, err = GetLocalTimeByCode(Filter{Code: 260}, t1)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(localTime)

	t.Log(GetLocalTimeByCode(Filter{Code: 1, CountryCode: "US"}, "2020-08-18 21:35:00"))
	t.Log(GetLocalTimeByCode(Filter{Code: 1, CountryCode: "CA"}, "2020-08-18 21:35:00"))
	t.Log(indexes[1])

	var date, _ = GetLocalTimeByCode(Filter{Code: 1, CountryCode: "US"}, time.Now().Format("2006-01-02 15:04:05"))
	t.Log(time.Unix(date, 0).Format("2006-01-02 15:04:05"))
}
