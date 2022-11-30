package integration

import (
	"net/http"
	"testing"

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
		res: response{code: http.StatusOK, content: content1},
	})
	sendRequestExpectResponse(t, requestResponse{
		req: request{method: http.MethodGet, url: "http://localhost:23443/other/test", body: http.NoBody},
		res: response{code: http.StatusOK, content: content2},
	})
}
