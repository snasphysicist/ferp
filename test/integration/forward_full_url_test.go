package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestForwardsQueryParametersOnDownstreamRequest(t *testing.T) {
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: echoQueryParameters()},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	expect, err := json.Marshal(map[string][]string{
		"foo": {"bar"},
		"oof": {"rab", "ferp"},
	})
	if err != nil {
		panic(err)
	}

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodGet,
			url:    proxyURL(p, "test") + "?foo=bar&oof=rab&oof=ferp",
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusOK,
			content: stringMatch{expect: string(expect)},
			headers: checkNoHeaders{},
		},
	})
}
