package wrappers

import "errors"

type DoubleValue float64

type FloatValue float32

type Int64Value int64

type UInt64Value uint64

type Int32Value int32

type UInt32Value uint32

type BoolValue bool

type StringValue string

func (*StringValue) ImplementsGraphQLType(name string) bool {
	return name == "String"
}

func (t *StringValue) UnmarshalGraphQL(input interface{}) error {
	if v, ok := input.(string); !ok {
		return errors.New("input value is not string")
	} else {
		*t = StringValue(v)
		return nil
	}
}

type BytesValue []byte

func WrapFloat64(f *float64) *DoubleValue {
	if f == nil {
		return nil
	}
	ret := DoubleValue(*f)
	return &ret
}
