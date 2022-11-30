package integration

import (
	"net/http"
	"testing"
)

func TestDoesNotForwardToUnconfiguredPath(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, code: 200, content: content},
		{path: "/teest", method: http.MethodGet, code: 200, content: content},
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
