package integration

import (
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
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
			headers: checkNoHeaders{},
		},
	})
}

func TestDoesNotCopyHeadersThatProxiesShouldDropToDownstreamRequest(t *testing.T) {
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: echoHeaders()},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	headers := http.Header{
		"Connection":     []string{"bar"},
		"Keep-Alive":     []string{"rab"},
		"Content-Length": []string{"1025"},
		"Close":          []string{"foo"},
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
			content: ensureDoesNotContainJSONSerialisedHeaders{expect: headers},
			headers: checkNoHeaders{},
		},
	})
}

func TestAddsForwardedHeaderWhenNoneIncoming(t *testing.T) {
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: echoHeaders()},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method:  http.MethodGet,
			url:     proxyURL(p, "test"),
			body:    http.NoBody,
			headers: http.Header{},
		},
		res: response{
			code: http.StatusOK,
			content: ensureJSONSerialisedForwardedHeaderMatching{
				matching: []string{
					`^by=ferp;for=(127\.0\.0\.1|\[::1\]):\d+;host=localhost:\d+;proto=HTTP/1\.1$`,
				},
			},
			headers: checkNoHeaders{},
		},
	})
}

func TestAppendsToForwardedHeaderWhenOneIncomingTest(t *testing.T) {
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, rg: echoHeaders()},
	}}

	p, f := startMocksAndProxy(t, []mock{m})
	defer f()

	sendRequestExpectResponse(t, requestResponse{
		req: request{
			method: http.MethodGet,
			url:    proxyURL(p, "test"),
			body:   http.NoBody,
			headers: http.Header{
				"Forwarded": []string{
					"for=192.0.2.43;proto=https;by=203.0.113.43;host=snas.pw",
				},
			},
		},
		res: response{
			code: http.StatusOK,
			content: ensureJSONSerialisedForwardedHeaderMatching{
				matching: []string{
					"^for=192\\.0\\.2\\.43;proto=https;by=203\\.0\\.113\\.43;host=snas\\.pw$",
					`^by=ferp;for=(127\.0\.0\.1|\[::1\]):\d+;host=localhost:\d+;proto=HTTP/1\.1$`,
				},
			},
			headers: checkNoHeaders{},
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

// ensureDoesNotContainJSONSerialisedHeaders fails the test if
// the body of the response does not contain serialised HTTP headers,
// and if those headers do contain any of the keys provided
type ensureDoesNotContainJSONSerialisedHeaders struct {
	expect http.Header
}

// Check implements contentMatcher for ensureContainsJSONSerialisedHeaders
func (m ensureDoesNotContainJSONSerialisedHeaders) Check(t *testing.T, b []byte) {
	var actual http.Header
	err := json.Unmarshal(b, &actual)
	if err != nil {
		t.Errorf("Failed to deserialise body to headers: %s", err)
		return
	}
	for k := range m.expect {
		if _, ok := actual[k]; ok {
			t.Errorf("Header key '%s' present in headers in body", k)
		}
	}
}

type ensureJSONSerialisedForwardedHeaderMatching struct {
	matching []string
}

func (m ensureJSONSerialisedForwardedHeaderMatching) Check(t *testing.T, b []byte) {
	var actual http.Header
	err := json.Unmarshal(b, &actual)
	if err != nil {
		t.Errorf("Failed to deserialise body to headers: %s", err)
		return
	}
	actualForwarded := actual["Forwarded"]
	if actualForwarded == nil {
		t.Error("Forwarded header was nil")
		return
	}
	if len(actualForwarded) != len(m.matching) {
		t.Errorf(
			"Forwarded header has %d values, expected %d",
			len(actualForwarded),
			len(m.matching),
		)
		return
	}
	for i, m := range m.matching {
		re := regexp.MustCompile(m)
		if !re.MatchString(actualForwarded[i]) {
			t.Errorf(
				"Header %d value '%s' did not match regex '%s'",
				i,
				actualForwarded[i],
				m,
			)
		}
	}
}
