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

	defer startMocksAndProxy(t, []mock{m})()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodGet,
			url:    "http://localhost:23443/prefixed/teest",
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusNotFound,
			content: "404 page not found\n",
		},
	})
}

func TestDoesNotForwardToUnconfiguredMethod(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: setResponse(200, content)},
	}}

	defer startMocksAndProxy(t, []mock{m})()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodPut,
			url:    "http://localhost:23443/prefixed/test",
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusMethodNotAllowed,
			content: "",
		},
	})
}
