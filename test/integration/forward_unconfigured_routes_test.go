package integration

import (
	"net/http"
	"testing"
	"time"

	"snas.pw/ferp/v2/cmd"
	"snas.pw/ferp/v2/configuration"
)

func TestDoesNotForwardToUnconfiguredPath(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPort, routes: []route{
		{path: "/test", method: http.MethodGet, code: 200, content: content},
	}}
	shutdown := m.start()
	defer shutdown()

	c, err := configuration.Load(mustFindFile("forward_test.yaml", "."))
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	defer close(stop)
	go cmd.Serve(c, stop)

	req, err := http.NewRequest(http.MethodGet, "http://localhost:23443/teest", http.NoBody)
	if err != nil {
		panic(err)
	}
	res := doUntilResponse(req, 11, time.Millisecond)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Response code was %d, but expected not found", res.StatusCode)
	}
}

func TestDoesNotForwardToUnconfiguredMethod(t *testing.T) {
	content := "Reached the test route"
	m := mock{t: t, port: mockPort, routes: []route{
		{path: "/test", method: http.MethodGet, code: 200, content: content},
	}}
	shutdown := m.start()
	defer shutdown()

	c, err := configuration.Load(mustFindFile("forward_test.yaml", "."))
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	defer close(stop)
	go cmd.Serve(c, stop)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:23443/test", http.NoBody)
	if err != nil {
		panic(err)
	}
	res := doUntilResponse(req, 11, time.Millisecond)
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Response code was %d, but expected method not allowed", res.StatusCode)
	}
}
