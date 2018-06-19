package timestamp

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes"
	"time"
	"fmt"
	"database/sql/driver"
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
	switch input := input.(type) {
	case *Timestamp:
		t = input
		return nil
	case Timestamp:
		t = &input
		return nil
	case time.Time:
		t = TimeToTimestamp(input)
		return nil
	case string:
		if val, err := time.Parse(time.RFC3339, input); err != nil {
			return err
		} else {
			t = TimeToTimestamp(val)
			return nil
		}
	case int:
		t = TimeToTimestamp(time.Unix(int64(input), 0))
		return nil
	case float64:
		t = TimeToTimestamp(time.Unix(int64(input), 0))
		return nil
	default:
		return fmt.Errorf("wrong type")
	}
}
