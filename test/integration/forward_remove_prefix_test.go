package integration

import (
	"net/http"
	"testing"
)

func TestForwardsWithPrefixRemovedWhenConfigured(t *testing.T) {
	content := "Reached a prefixed test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodPost, rg: setResponse(200, content)},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodPost,
			url:    proxyURL(p, "prefixed/test"),
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusOK,
			content: stringMatch{expect: content},
		},
	})
}
