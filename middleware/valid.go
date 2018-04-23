package middleware

import (
	"context"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
)

// AuthError represents an authorization error.
type ValidError struct {
	errMsg string
}

// StatusCode is an implementation of the StatusCoder interface in go-kit/http.
func (valid ValidError) StatusCode() int {
	return http.StatusNotAcceptable
}

// Error is an implementation of the Error interface.
func (valid ValidError) Error() string {
	// return http.StatusText(http.StatusNotAcceptable)
	return valid.errMsg
}

// Headers is an implementation of the Headerer interface in go-kit/http.
func (e ValidError) Headers() http.Header {
	return http.Header{
		"Content-Type":           []string{"text/plain; charset=utf-8"},
		"X-Content-Type-Options": []string{"nosniff"},
	}
}

func ValidMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			govalidator.CustomTypeTagMap.Set("uidarr", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
				if len(i.([]uint64)) > 0 {
					return true
				}
				return false
			}))
			result, err := govalidator.ValidateStruct(request)
			if err != nil {
				// println("error: " + err.Error())
				return nil, ValidError{err.Error()}
			}
			if !result {
				return nil, ValidError{err.Error()}
			}

			return next(ctx, request)
		}
	}
}
