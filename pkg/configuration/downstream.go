package configuration

import (
	"fmt"

	"github.com/snasphysicist/ferp/v2/pkg/functional"
)

// populateDownstreams finds downstreams for all incomings in the configuration,
// returning an error if some cannot be found
func populateDownstreams(c Configuration) (Configuration, error) {
	isi, isErr := findDownstreams(c.Downstreams, c.HTTP.Incoming)
	c.HTTP.Incoming = isi
	si, sErr := findDownstreams(c.Downstreams, c.HTTPS.Incoming)
	c.HTTPS.Incoming = si
	return c, joinNonNilErrors([]error{isErr, sErr}, ", ", "%s")
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
	err := joinNonNilErrors(errs, ", ", "invalid downstreams: %s")
	return iswd, err
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
