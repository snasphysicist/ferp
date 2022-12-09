package url

import (
	"net/url"
	"strings"
	"testing"

	"github.com/snasphysicist/ferp/v2/pkg/log"
	"github.com/snasphysicist/ferp/v2/pkg/mapper"
)

func TestChangesHostAndPortToTargetOnes(t *testing.T) {
	_, _ = log.Initialise()

	u := url.URL{Host: "something-else:1089"}
	b := BaseURL{Host: "target", Port: 8082}
	m := mapper.Passthrough{}
	s := Rewrite(u, b, m.Map)
	expect := "//target:8082"
	if s != expect {
		t.Errorf("Rewrote to '%s' not '%s' from '%#v' using '%#v' & '%#v'",
			s, expect, u, b, m)
	}
}

func TestUsesOnlyBasePathWhenOriginalPathIsRoot(t *testing.T) {
	_, _ = log.Initialise()

	u := url.URL{Host: "something-else:1089", Path: "/"}
	b := BaseURL{Host: "target", Port: 8082, Path: "/foo/bar/"}
	m := mapper.Passthrough{}
	s := Rewrite(u, b, m.Map)
	expect := "//target:8082/foo/bar/"
	if s != expect {
		t.Errorf("Rewrote to '%s' not '%s' from '%#v' using '%#v' & '%#v'",
			s, expect, u, b, m)
	}
}

func TestRewritesToRootWhenOriginalIsRootBaseIsEmpty(t *testing.T) {
	_, _ = log.Initialise()

	u := url.URL{Host: "something-else:1089", Path: "/"}
	b := BaseURL{Host: "target", Port: 8082, Path: ""}
	m := mapper.Passthrough{}
	s := Rewrite(u, b, m.Map)
	expect := "//target:8082/"
	if s != expect {
		t.Errorf("Rewrote to '%s' not '%s' from '%#v' using '%#v' & '%#v'",
			s, expect, u, b, m)
	}
}

func TestRewritesToRootWhenBaseIsRootOriginalIsEmpty(t *testing.T) {
	_, _ = log.Initialise()

	u := url.URL{Host: "something-else:1089", Path: ""}
	b := BaseURL{Host: "target", Port: 8082, Path: "/"}
	m := mapper.Passthrough{}
	s := Rewrite(u, b, m.Map)
	expect := "//target:8082/"
	if s != expect {
		t.Errorf("Rewrote to '%s' not '%s' from '%#v' using '%#v' & '%#v'",
			s, expect, u, b, m)
	}
}

func TestRewritesToEmptyPathWhenBothEmpty(t *testing.T) {
	_, _ = log.Initialise()

	u := url.URL{Host: "something-else:1089", Path: ""}
	b := BaseURL{Host: "target", Port: 8082, Path: ""}
	m := mapper.Passthrough{}
	s := Rewrite(u, b, m.Map)
	expect := "//target:8082"
	if s != expect {
		t.Errorf("Rewrote to '%s' not '%s' from '%#v' using '%#v' & '%#v'",
			s, expect, u, b, m)
	}
}

func TestJoinsOriginalAndBasePathsWithSingleSlash(t *testing.T) {
	_, _ = log.Initialise()

	for _, original := range []string{"foo", "/foo", "foo/", "/foo/"} {
		for _, base := range []string{"bar", "/bar", "bar/", "/bar/"} {
			u := url.URL{Host: "something-else:1089", Path: original}
			b := BaseURL{Host: "target", Port: 8082, Path: base}
			m := mapper.Passthrough{}
			s := Rewrite(u, b, m.Map)
			expectBeginning := "//target:8082/bar/foo"
			if !strings.HasPrefix(s, expectBeginning) {
				t.Errorf(
					"Rewrote to '%s' which does not begin with '%s'"+
						" from '%#v' using '%#v' & '%#v'",
					s, expectBeginning, u, b, m)
			}
		}
	}
}

func TestRetainsOriginalQueryParameters(t *testing.T) {
	query := "foo=bar&baz=rab"
	u := url.URL{Host: "something-else:1089", Path: "/foo/", RawQuery: query}
	b := BaseURL{Host: "target", Port: 8082, Path: "/bar/"}
	m := mapper.Passthrough{}
	s := Rewrite(u, b, m.Map)
	if !strings.HasSuffix(s, query) {
		t.Errorf(
			"Rewrote to '%s' which does not end with '%s'"+
				" from '%#v' using '%#v' & '%#v'",
			s, query, u, b, m)
	}
}
