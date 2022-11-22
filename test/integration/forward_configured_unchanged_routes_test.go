package integration

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"snas.pw/ferp/v2/cmd"
	"snas.pw/ferp/v2/configuration"
	"snas.pw/ferp/v2/functional"
)

func TestForwardsOnConfiguredExactRouteAndMethod(t *testing.T) {
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
	m1 := mock{t: t, port: mockPort, routes: []route{
		{path: "/test", method: http.MethodGet, code: 200, content: content1},
	}}
	shutdown := m1.start()
	defer shutdown()

	content2 := "Reached the other test route"
	m2 := mock{t: t, port: otherMockPort, routes: []route{
		{path: "/other/test", method: http.MethodGet, code: 200, content: content2},
	}}
	shutdown = m2.start()
	defer shutdown()

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

const mockPort = 34543

const otherMockPort = 34555

func mustFindFile(nameWithExtension string, startIn string) string {
	current, err := filepath.Abs(startIn)
	if err != nil {
		panic(err)
	}
	all := make([]string, 0)
	err = filepath.Walk(current, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		all = append(all, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	matches := functional.Filter(all,
		func(s string) bool { return strings.HasSuffix(s, nameWithExtension) })
	if len(matches) != 1 {
		return mustFindFile(nameWithExtension, parentOf(current))
	}
	return matches[0]
}

func parentOf(path string) string {
	if path == "/" {
		panic("top level directory / has no parent")
	}
	return filepath.Dir(path)
}

func doUntilResponse(r *http.Request, retries uint, backoff time.Duration) *http.Response {
	var i uint
	for i < retries {
		i++
		res, err := (&http.Client{}).Do(r)
		if err == nil {
			return res
		}
		time.Sleep(backoff)
		backoff = backoff * 2
	}
	panic("no sucessful request after retries")
}
