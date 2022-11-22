package integration

import (
	"io"
	"net/http"
	"testing"
	"time"

	"snas.pw/ferp/v2/cmd"
	"snas.pw/ferp/v2/configuration"
)

func TestForwardsWithPrefixRemovedWhenConfigured(t *testing.T) {
	content := "Reached a prefixed test route"
	m := mock{t: t, port: mockPort, routes: []route{
		{path: "/test", method: http.MethodPost, code: 200, content: content},
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

	req, err := http.NewRequest(http.MethodPost, "http://localhost:23443/prefixed/test", http.NoBody)
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
