package url

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// Rewrite adapts the provided URL to the new target,
// mapping the original path to the new one and joining
// this with the base path, taking care of e.g. repeated slashes
func Rewrite(u url.URL, target BaseURL, pm PathRewriter) string {
	u.Scheme = target.Protocol
	u.Host = fmt.Sprintf("%s:%d", target.Host, target.Port)
	mapped := pm(u.Path)
	u.Path = assembleFullPath(target.Path, mapped)
	return u.String()
}

// PathRewriter maps between the incoming path and the outgoing path
type PathRewriter func(string) string

// assembleFullPath combines the base and suffix into a single path
// ensuring they are joined by a single slash
func assembleFullPath(base string, suffix string) string {
	joined := joinPreservingTrailingSlash(base, suffix)
	if joined == "" {
		return joined
	}
	return fmt.Sprintf("/%s", strings.TrimLeft(joined, "/"))
}

// joinPreservingTrailingSlash wraps path.Join, preserving the trailing slash if any
func joinPreservingTrailingSlash(base string, suffix string) string {
	joined := path.Join(base, strings.TrimLeft(suffix, "/"))
	if strings.HasSuffix(suffix, "/") {
		return fmt.Sprintf("%s/", joined)
	}
	return joined
}

// BaseURL represents the "base" part of a URL
// to which a subpath (suffix) can be added
type BaseURL struct {
	Protocol string
	Host     string
	Port     uint16
	Path     string
}
