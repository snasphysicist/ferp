package integration

import (
	"net/http"
	"testing"
)

func TestGracefullyFailsWhenDownstreamIsUnavailable(t *testing.T) {
	startMocksAndProxy(t, []mock{})

	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: "http://localhost:23443/test", body: http.NoBody},
		res: response{code: http.StatusInternalServerError, content: "500: something went wrong"},
	})
}
