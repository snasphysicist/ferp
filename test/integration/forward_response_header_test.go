package integration

import (
	"net/http"
	"reflect"
	"testing"
)

func TestForwardAllResponseHeadersExceptDroppedOnes(t *testing.T) {
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: withSetHeaders()},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	// The bug this catches happens intermittently, due to randomised map order
	for i := 0; i < 100; i++ {
		sendRequestExpectResponse(t, requestResponse{
			req: request{
				method:  http.MethodGet,
				url:     proxyURL(p, "test"),
				body:    http.NoBody,
				headers: make(http.Header),
			},
			res: response{
				code:    http.StatusOK,
				content: stringMatch{expect: ""},
				headers: checkContainsSetHeaders{},
			},
		})
	}
}

// withSetHeaders generates a response with 200 status, empty body, and some known headers.
// One of the headers should not be forwarded along by the proxy, one of them should.
func withSetHeaders() responseGenerator {
	return func(_ *http.Request) responseSpecification {
		return responseSpecification{status: 200, body: "", headers: http.Header{
			"Connection":   []string{"Shouldn't-Forward"},
			"Content-Type": []string{"application/json"},
		}}
	}
}

// checkContainsSetHeaders fails the test if the response does not contain
// the headers that should be set on the response from withSetHeaders
type checkContainsSetHeaders struct{}

// Check implements headerMatcher for checkContainsSetHeaders, see struct for behaviour
func (checkContainsSetHeaders) Check(t *testing.T, h http.Header) {
	if _, ok := h["Connection"]; ok {
		t.Fatalf("Connection header was forwarded in %#v", h)
	}
	vs, ok := h["Content-Type"]
	if !ok {
		t.Fatalf("Content-Type header was not forwarded in %#v", h)
	}
	if !reflect.DeepEqual(vs, []string{"application/json"}) {
		t.Fatalf("Content-Type values %v, expected one value application/json", vs)
	}
}
