package configuration

import (
	"fmt"
	"strings"

	"github.com/snasphysicist/ferp/v2/pkg/functional"
)

// populateDownstreams finds downstreams for all incomings in the configuration,
// returning an error if some cannot be found
func populateDownstreams(c Configuration) (Configuration, error) {
	isi, err := findDownstreams(c.Downstreams, c.HTTP.Incoming)
	c.HTTP.Incoming = isi
	errs := []error{err}
	errs = functional.Filter(errs, func(e error) bool { return e != nil })
	errStr := functional.Map(errs, func(e error) string { return e.Error() })
	if len(errs) != 0 {
		return c, fmt.Errorf("%s", strings.Join(errStr, ", "))
	}
	return c, nil
}

// findDownstreams finds downstreams for all provided incomings, returning
// an error describing all those that could not be found (if any)
func findDownstreams(d []Downstream, is []Incoming) ([]Incoming, error) {
	iswd := make([]Incoming, 0)
	errs := make([]error, 0)
	for _, i := range is {
		d, err := findDownstream(d, i)
		errs = append(errs, err)
		iswd = append(iswd, Incoming{
			Path:          i.Path,
			Methods:       i.Methods,
			MethodRouters: i.MethodRouters,
			Target:        i.Target,
			Downstream:    d,
		})
	}
	errStr := functional.Map(
		functional.Filter(errs, func(e error) bool { return e != nil }),
		func(e error) string { return e.Error() })
	if len(errs) > 0 {
		return is,
			fmt.Errorf("invalid downstreams: %s", strings.Join(errStr, ", "))
	}
	return iswd, nil
}

// findDownstream attempts to find the downstream in the configuration for the given incoming route
func findDownstream(d []Downstream, i Incoming) (Downstream, error) {
	matches := functional.Filter(d, func(d Downstream) bool { return d.Target == i.Target })
	if len(matches) != 1 {
		return Downstream{}, fmt.Errorf(
			"invalid downstream target %s from incoming %+v", i.Target, i)
	}
	return matches[0], nil
}
