package integration

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/snasphysicist/ferp/v2/pkg/functional"
)

// mustFindFile searches recursively, starting from startIn and moving through its parents,
// for the file with the given name.extension.
// Returns the absolute path when found, panics if the file cannot be found.
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

// parentOf tries to find the parent of the provided path, panicing if it looks like none exists.
func parentOf(path string) string {
	lastPathCharacter := path[len(path)-1]
	// see filepath.Dir documentation to understand why this is true
	isTopMostDirectory := os.IsPathSeparator(lastPathCharacter)
	if isTopMostDirectory {
		panic("top level directory / has no parent")
	}
	return filepath.Dir(path)
}

// doUntilResponse repeatedly sends the request, with exponentially increasing backoff starting
// from backoff, up to retries times, until no error is encountered.
// panics if no there is no successful request after retries attempts.
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

// sendRequestExpectResponse sends a request with the given method, url and body
// and fails the test if the response does not have the given status code and content,
// or if anything at all goes wrong in the request-response cycle.
func sendRequestExpectResponse(t *testing.T, rr requestResponse) {
	t.Helper()
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
