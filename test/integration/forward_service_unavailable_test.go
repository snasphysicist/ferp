package integration

import (
	"net/http"
	"testing"
	"time"

	"snas.pw/ferp/v2/cmd"
	"snas.pw/ferp/v2/configuration"
)

func TestGracefullyFailsWhenDownstreamIsUnavailable(t *testing.T) {
	c, err := configuration.Load(mustFindFile("forward_test.yaml", "."))
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	defer close(stop)
	go cmd.Serve(c, stop)

	req, err := http.NewRequest(http.MethodGet, "http://localhost:23443/test", http.NoBody)
	if err != nil {
		panic(err)
	}
	res := doUntilResponse(req, 11, time.Millisecond)
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Response code was %d, but expected internal error", res.StatusCode)
	}
}
