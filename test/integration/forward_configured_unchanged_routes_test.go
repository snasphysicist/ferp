package integration

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/snasphysicist/ferp/v2/command"
	"github.com/snasphysicist/ferp/v2/pkg/configuration"
	"github.com/snasphysicist/ferp/v2/pkg/log"
)

func TestForwardsOnConfiguredExactRouteAndMethod(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodGet, code: 200, content: content},
	}}
	shutdown := m.start()
	defer shutdown()

	_, _ = log.Initialise()
	c, err := configuration.Load(mustFindFile("test.yaml", "."))
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	defer close(stop)
	go command.Serve(c, stop)

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
	shutdown := m1.start()
	defer shutdown()

	content2 := "Reached the other test route"
	m2 := mock{t: t, port: mockPorts()[1], routes: []route{
		{path: "/other/test", method: http.MethodGet, code: 200, content: content2},
	}}
	shutdown = m2.start()
	defer shutdown()

	c, err := configuration.Load(mustFindFile("test.yaml", "."))
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	defer close(stop)
	go command.Serve(c, stop)

	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: "http://localhost:23443/test", body: http.NoBody},
		res: response{code: http.StatusOK, content: content1},
	})
	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: "http://localhost:23443/other/test", body: http.NoBody},
		res: response{code: http.StatusOK, content: content2},
	})
}

// sendRequestExpectResponse sends a request with the given method, url and body
// and fails the test if the response does not have the given status code and content,
// or if anything at all goes wrong in the request-response cycle.
func sendRequestExpectResponse(t *testing.T, rr requestResponse) {
	req, err := http.NewRequest(rr.req.method, rr.req.url, rr.req.body)
	if err != nil {
		t.Errorf("Failed to construct request: %s", err)
		return
	}
	res := doUntilResponse(req, 11, time.Millisecond)
	if res.StatusCode != rr.res.code {
		t.Errorf("Request had status %d, expected %d", res.StatusCode, rr.res.code)
		return
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %s", err)
		return
	}
	if string(b) != rr.res.content {
		t.Errorf("Response content does not match expected: '%s' != '%s'",
			string(b), rr.res.content)
	}
}

// requestResponse stores a request to be sent during a test and the response expected
type requestResponse struct {
	req request
	res response
}

// request represents a request to be send during a test
type request struct {
	method string
	url    string
	body   io.ReadCloser
}

// response represents the expected state of a response to be returned during a test
type response struct {
	code    int
	content string
}
