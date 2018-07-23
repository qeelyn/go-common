package errors

import (
	"errors"
)

// graphql input is not except type
func GqlInputWrongType() error {
	return errors.New("GQL_NOT_EXCEPT_TYPE")
}
