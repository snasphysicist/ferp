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

	req, err := http.NewRequest(http.MethodGet, "http://localhost:23443/test", http.NoBody)
	if err != nil {
		panic(err)
	}
	res := doUntilResponse(req, 11, time.Millisecond)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Request failed with %d, but should have been forwarded", res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if string(b) != content {
		t.Errorf("Forwarded content %s != sent content %s", string(b), content)
	}
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

	req, err := http.NewRequest(http.MethodGet, "http://localhost:23443/test", http.NoBody)
	if err != nil {
		panic(err)
	}
	res := doUntilResponse(req, 11, time.Millisecond)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Request failed with %d, but should have been forwarded", res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if string(b) != content1 {
		t.Errorf("Forwarded content %s != sent content %s", string(b), content1)
	}

	req2, err := http.NewRequest(http.MethodGet, "http://localhost:23443/other/test", http.NoBody)
	if err != nil {
		panic(err)
	}
	res2 := doUntilResponse(req2, 11, time.Millisecond)
	if res2.StatusCode != http.StatusOK {
		t.Errorf("Request failed with %d, but should have been forwarded", res2.StatusCode)
	}
	b, err = io.ReadAll(res2.Body)
	if err != nil {
		panic(err)
	}
	if string(b) != content2 {
		t.Errorf("Forwarded content %s != sent content %s", string(b), content2)
	}
}
