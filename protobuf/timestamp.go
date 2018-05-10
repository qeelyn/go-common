package protobuf

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes"
	"time"
	"fmt"
)

type Timestamp timestamp.Timestamp

func (g *Timestamp) Scan(src interface{}) error {
	t, _ := ptypes.TimestampProto(src.(time.Time))
	g.Seconds = t.Seconds
	g.Nanos = t.Nanos
	return nil
}

func (g *Timestamp) ToTime() time.Time {
	val, err := ptypes.Timestamp((*timestamp.Timestamp)(g))
	if err != nil {
		return time.Time{}
	}
	return val
}

func (g *Timestamp) ParseFrom(t time.Time) error {
	tmp, err := ptypes.TimestampProto(t)
	if err != nil {
		return err
	}
	g.Seconds = tmp.Seconds
	g.Nanos = tmp.Nanos
	return nil
}

func (g *Timestamp) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", g.ToTime().Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}