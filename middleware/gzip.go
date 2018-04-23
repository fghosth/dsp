package middleware

import (
	"compress/gzip"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

// AuthError represents an authorization error.
type GzipError struct {
	errMsg string
}

// StatusCode is an implementation of the StatusCoder interface in go-kit/http.
func (gzip GzipError) StatusCode() int {
	return http.StatusGone
}

// Error is an implementation of the Error interface.
func (gzip GzipError) Error() string {
	// return http.StatusText(http.StatusNotAcceptable)
	return gzip.errMsg
}

// Headers is an implementation of the Headerer interface in go-kit/http.
func (e GzipError) Headers() http.Header {
	return http.Header{
		"Content-Type":           []string{"text/plain; charset=utf-8"},
		"X-Content-Type-Options": []string{"nosniff"},
	}
}

func GzipMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			var encoding string
			r, ok := request.(*http.Request)
			if ok {
				encoding = r.Header.Get("Content-Encoding")
			}
			if encoding == "gzip" {
				resp, err := gzip.NewReader(r.Body)
				if err != nil {
					return nil, GzipError{"body不是gzip格式"}
				}
				body, _ := ioutil.ReadAll(resp)
				ctx = context.WithValue(ctx, "body", body) //线程安全，并发性能受到影响
			}
			return next(ctx, request)
		}
	}
}
