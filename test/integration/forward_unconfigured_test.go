package integration

import (
	"net/http"
	"testing"
)

func TestDoesNotForwardToUnconfiguredPath(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: setResponse(200, content)},
		{path: "/teest", method: http.MethodGet, rg: setResponse(200, content)},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodGet,
			url:    proxyURL(p, "prefixed/teest"),
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusNotFound,
			content: stringMatch{expect: "404 page not found\n"},
			headers: checkNoHeaders{},
		},
	})
}

func TestDoesNotForwardToUnconfiguredMethod(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: setResponse(200, content)},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodPut,
			url:    proxyURL(p, "prefixed/test"),
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusMethodNotAllowed,
			content: stringMatch{expect: ""},
			headers: checkNoHeaders{},
		},
	})
}
