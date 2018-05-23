package protobuf_test

import (
	"testing"
	"time"
	"github.com/qeelyn/go-common/protobuf"
	"fmt"
)

func TestTimeToTimestamp(t *testing.T) {
	exp := time.Date(2018,5,16,9,50,30,0,time.UTC)
	exp = time.Now()
	tsmp := protobuf.TimeToTimestamp(exp)

	if tsmp.Seconds != exp.Unix() {
		t.Error("t to tsmp error!")
	}
}

func TestTimestampToTime(t *testing.T) {
	exp := time.Date(2018,5,16,9,50,30,0,time.UTC)
	exp = time.Now()
	tsmp := protobuf.TimeToTimestamp(exp)

	tm,err := protobuf.TimestampToTime(tsmp)
	if err != nil {
		t.Fatal(err)
	}
	if tm.IsZero() {
		t.Fatal("convert error")
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	exp := time.Date(2018,5,16,9,50,30,0,time.Local)
	tsmp := protobuf.TimeToTimestamp(exp)
	s,_:= tsmp.MarshalJSON()
	act := string(s)
	want := fmt.Sprintf("\"%s\"",exp.Local().Format(time.RFC3339))
	println(act)
	println(want)
	if act != want {
		t.Fatal()
	}
}
