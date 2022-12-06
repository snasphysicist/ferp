package integration

import (
	"net/http"
	"testing"
)

func TestGracefullyFailsWhenDownstreamIsUnavailable(t *testing.T) {
	p, f := startMocksAndProxy(t, []mock{})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: proxyURL(p, "test"), body: http.NoBody},
		res: response{code: http.StatusInternalServerError, content: "500: something went wrong"},
	})
}
