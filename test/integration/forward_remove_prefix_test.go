package integration

import (
	"net/http"
	"testing"

	"github.com/snasphysicist/ferp/v2/command"
	"github.com/snasphysicist/ferp/v2/pkg/configuration"
	"github.com/snasphysicist/ferp/v2/pkg/log"
)

func TestForwardsWithPrefixRemovedWhenConfigured(t *testing.T) {
	content := "Reached a prefixed test route"
	m := mock{t: t, port: mockPorts()[0], routes: []route{
		{path: "/test", method: http.MethodPost, code: 200, content: content},
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
		req: request{
			method: http.MethodPost,
			url:    "http://localhost:23443/prefixed/test",
			body:   http.NoBody,
		},
		res: response{
			code:    http.StatusOK,
			content: content,
		},
	})
}
