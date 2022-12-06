package integration

import (
	"net/http"
	"testing"
)

func TestForwardsOnConfiguredExactRouteAndMethod(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: setResponse(200, content)},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: proxyURL(p, "test"), body: http.NoBody},
		res: response{code: http.StatusOK, content: content},
	})
}

func TestForwardsPathsToCorrespondingDownstreams(t *testing.T) {
	content1 := "Reached the test route"
	m1 := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: setResponse(200, content1)},
	}}

	content2 := "Reached the other test route"
	m2 := mock{t: t, port: mockPorts()[1], routes: []route{
		{path: "/other/test", method: http.MethodGet, rg: setResponse(200, content2)},
	}}

	p, f := startMocksAndProxy(t, []mock{m1, m2})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: proxyURL(p, "test"), body: http.NoBody},
		res: response{code: http.StatusOK, content: content1},
	})
	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: proxyURL(p, "other/test"), body: http.NoBody},
		res: response{code: http.StatusOK, content: content2},
	})
}
