package integration

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestCopiesHeadersToDownstreamRequest(t *testing.T) {
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: echoHeaders()},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	headers := http.Header{
		"Foo": []string{"bar"},
		"Oof": []string{"rab", "ferp"},
	}

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method:  http.MethodGet,
			url:     proxyURL(p, "test"),
			body:    http.NoBody,
			headers: headers,
		},
		res: response{
			code:    http.StatusOK,
			content: ensureContainsJSONSerialisedHeaders{expect: headers},
		},
	})

}

// ensureContainsJSONSerialisedHeaders fails the test if the body of the response
// does not contain serialised HTTP headers, and if those headers do not
// contain all the key value pairs provided
type ensureContainsJSONSerialisedHeaders struct {
	expect http.Header
}

// Check implements contentMatcher for ensureContainsJSONSerialisedHeaders
func (m ensureContainsJSONSerialisedHeaders) Check(t *testing.T, b []byte) {
	var actual http.Header
	err := json.Unmarshal(b, &actual)
	if err != nil {
		t.Errorf("Failed to deserialise body to headers: %s", err)
		return
	}
	for k, vs := range m.expect {
		if _, ok := actual[k]; !ok {
			t.Errorf("Header key '%s' missing from headers in body", k)
		}
		if !reflect.DeepEqual(vs, actual[k]) {
			t.Errorf("Header value for key '%s' got '%+v' expected '%+v'",
				k, actual[k], vs)
		}
	}
}
