package timestamp

import (
	"database/sql/driver"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/qeelyn/go-common/protobuf/errors"
	"time"
)

type Timestamp timestamp.Timestamp

func (g *Timestamp) ToTime() (time.Time, error) {
	return TimestampToTime(g)
}

// database scan
func (g *Timestamp) Scan(src interface{}) error {
	t, _ := ptypes.TimestampProto(src.(time.Time))
	g.Seconds = t.Seconds
	g.Nanos = t.Nanos
	return nil
}

// database Valuer Interface
func (g Timestamp) Value() (driver.Value, error) {
	return g.ToTime()
}

// json support
func (g Timestamp) MarshalJSON() ([]byte, error) {
	t, err := TimestampToTime(&g)
	if err != nil {
		return nil, err
	}
	var stamp = "\"" + t.Local().Format(time.RFC3339) + "\""
	return []byte(stamp), nil
}

func (Timestamp) ImplementsGraphQLType(name string) bool {
	return name == "Time"
}

func (t *Timestamp) UnmarshalGraphQL(input interface{}) error {
	switch in := input.(type) {
	case *Timestamp:
		t = in
		return nil
	case Timestamp:
		t = &in
		return nil
	case time.Time:
		t = TimeToTimestamp(in)
		return nil
	case string:
		if val, err := time.Parse(time.RFC3339, in); err != nil {
			return err
		} else {
			t = TimeToTimestamp(val)
			return nil
		}
	case int:
		t = TimeToTimestamp(time.Unix(int64(in), 0))
		return nil
	case float64:
		t = TimeToTimestamp(time.Unix(int64(in), 0))
		return nil
	default:
		return errors.GqlInputWrongType()
	}
}

// proto message
func (t *Timestamp) Reset() {
	*t = Timestamp{}
}

func (*Timestamp) ProtoMessage() {}

func (t *Timestamp) String() string {
	return proto.CompactTextString(t)
}
