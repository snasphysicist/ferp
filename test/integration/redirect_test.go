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
			headers: checkNoHeaders{},
		},
	})
}

func TestConfiguredRedirectReturnsLocationHeaderWithTargetURL(t *testing.T) {
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
			headers: checkLocationHeader{content: "/you-are-redirected"},
		},
	})
}

// checkNothing is a contentMatcher that never fails the test
type checkNothing struct{}

// Check implements contentMatcher for checkNothing, see struct for behaviour
func (checkNothing) Check(t *testing.T, b []byte) {}

// checkLocationHeader is a headerMatcher which expects the Location header
// to contain precisely one value, which matches exactly the wrapped string
type checkLocationHeader struct {
	content string
}

// Check implements headerMatcher for checkLocationHeader, see struct for behaviour
func (c checkLocationHeader) Check(t *testing.T, h http.Header) {
	locations, ok := h["Location"]
	if !ok {
		t.Errorf("No location header in %#v", h)
		return
	}
	if len(locations) != 1 {
		t.Errorf("location header contains %d values, expected 1 (%#v)",
			len(locations), locations)
		return
	}
	if locations[0] != c.content {
		t.Errorf("location header value '%s', expected '%s'",
			locations[0], c.content)
	}
}
