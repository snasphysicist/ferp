package integration

import (
	"net/http"
	"testing"
)

func TestRedirectStatusReturnedOnConfiguredRedirectPath(t *testing.T) {
	p, f := startMocksAndProxy(t, []mock{})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodGet,
			url:    proxyURL(p, "redirect-me"),
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusFound,
			content: checkNothing{},
		},
	})
}

// checkNothing is a contentMatcher that never fails the test
type checkNothing struct{}

// Check implements contentMatcher for checkNothing, see struct for behaviour
func (checkNothing) Check(t *testing.T, b []byte) {}
