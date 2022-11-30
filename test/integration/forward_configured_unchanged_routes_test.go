package integration

import (
	"net/http"
	"testing"
)

func TestForwardsOnConfiguredExactRouteAndMethod(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, code: 200, content: content},
	}}

	defer startMocksAndProxy(t, []mock{m})()

	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: "http://localhost:23443/test", body: http.NoBody},
		res: response{code: http.StatusOK, content: content},
	})
}

func TestForwardsPathsToCorrespondingDownstreams(t *testing.T) {
	content1 := "Reached the test route"
	m1 := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, code: 200, content: content1},
	}}

	content2 := "Reached the other test route"
	m2 := mock{t: t, port: mockPorts()[1], routes: []route{
		{path: "/other/test", method: http.MethodGet, code: 200, content: content2},
	}}

	defer startMocksAndProxy(t, []mock{m1, m2})()

	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: "http://localhost:23443/test", body: http.NoBody},
		res: response{code: http.StatusOK, content: content1},
	})
	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: "http://localhost:23443/other/test", body: http.NoBody},
		res: response{code: http.StatusOK, content: content2},
	})
}
