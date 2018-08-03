package date

import (
	"fmt"
	"time"

	"github.com/qeelyn/go-common/protobuf/errors"
)

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

func (t *Date) TimeString() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("%d-%02d-%02d", t.Year, t.Month, t.Day)
}
