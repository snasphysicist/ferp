package integration

import (
	"net/http"
	"testing"
)

func TestForwardsWithPrefixRemovedWhenConfigured(t *testing.T) {
	content := "Reached a prefixed test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodPost, code: 200, content: content},
	}}

	defer startMocksAndProxy(t, []mock{m})()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodPost,
			url:    "http://localhost:23443/prefixed/test",
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusOK,
			content: content,
		},
	})
}
