package date

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"

	"github.com/qeelyn/go-common/protobuf/errors"
)

type Date struct {
	Year  int32 `protobuf:"varint,1,opt,name=year" json:"year,omitempty"`
	Month int32 `protobuf:"varint,2,opt,name=month" json:"month,omitempty"`
	Day   int32 `protobuf:"varint,3,opt,name=day" json:"day,omitempty"`
}

func (Date) ImplementsGraphQLType(name string) bool {
	return name == "Date"
}

func (t *Date) UnmarshalGraphQL(input interface{}) error {
	switch ts := input.(type) {
	case string:
		if val, err := time.Parse("2006-01-02", ts); err != nil {
			return err
		} else {
			t.Year = int32(val.Year())
			t.Month = int32(val.Month())
			t.Day = int32(val.Day())
			return nil
		}
	default:
		return errors.GqlInputWrongType()
	}
}

func (t Date) ToString() string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year, t.Month, t.Day)
}

// proto message
func (t *Date) Reset() {
	*t = Date{}
}

func (*Date) ProtoMessage() {}

func (t *Date) String() string {
	return proto.CompactTextString(t)
}
