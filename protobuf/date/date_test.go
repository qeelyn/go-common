package date_test

import (
	"github.com/qeelyn/go-common/protobuf/date"
	"testing"
)

func TestDate_ToString(t *testing.T) {
	d := &date.Date{
		Year:  2018,
		Month: 1,
		Day:   11,
	}
	if d.ToString() != "2018-01-11" {
		t.Error("tostring error")
	}
}

func TestDate_UnmarshalGraphQL(t *testing.T) {
	expect := "2018-01-11"
	d := &date.Date{}
	d.UnmarshalGraphQL(expect)
	if d.ToString() != expect {
		t.Error()
	}
}
