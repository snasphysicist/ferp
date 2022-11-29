package integration

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
