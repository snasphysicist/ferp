package configuration

import (
	"fmt"
	"strings"

	"github.com/snasphysicist/ferp/v2/pkg/configuration/router"
	"github.com/snasphysicist/ferp/v2/pkg/functional"
	"github.com/snasphysicist/ferp/v2/pkg/log"
)

// populateMethodRouters adds method routers for all configured redirects and forwarded routes
func populateMethodRouters(c Configuration) (Configuration, error) {
	ii, iiErr := populateMethodRoutersForIncomings(c.HTTP.Incoming)
	c.HTTP.Incoming = ii
	errs := []error{iiErr}
	errs = functional.Filter(errs, func(e error) bool { return e != nil })
	errStr := functional.Map(errs, func(e error) string { return e.Error() })
	if len(errs) > 0 {
		return c, fmt.Errorf("invalid methods: %s", strings.Join(errStr, ", "))
	}
	return c, nil
}

// populateMethodRoutersForIncomings populates the method routers fields in the provided
// incomings, returning an error summarising which (if any) had invalid methods
func populateMethodRoutersForIncomings(is []Incoming) ([]Incoming, error) {
	iswr := make([]Incoming, 0)
	errs := make([]error, 0)
	for _, i := range is {
		mrs, err := findMethodRouters(i.Methods)
		errs = append(errs, err)
		iswr = append(iswr, Incoming{
			Path:          i.Path,
			Methods:       i.Methods,
			MethodRouters: mrs,
			Target:        i.Target,
			Downstream:    i.Downstream,
		})
		log.Infof("For %+v, method routers %+v", i, mrs)
	}
	errs = functional.Filter(errs, func(e error) bool { return e != nil })
	errStr := functional.Map(errs, func(e error) string { return e.Error() })
	if len(errs) > 0 {
		return is, fmt.Errorf("invalid methods: %s", strings.Join(errStr, ", "))
	}
	return iswr, nil
}

// findMethodRouters finds method routers for all provided methods,
// returning an error describing for which methods (if any) routers could not be found
func findMethodRouters(ms []string) ([]router.MethodRouter, error) {
	errs := make([]error, 0)
	rs := make([]router.MethodRouter, 0)
	for _, m := range ms {
		r, err := router.RouterFor(m)
		rs = append(rs, r)
		errs = append(errs, err)
	}
	errs = functional.Filter(errs, func(e error) bool { return e != nil })
	errStr := functional.Map(errs, func(e error) string { return e.Error() })
	if len(errs) != 0 {
		return []router.MethodRouter{},
			fmt.Errorf("%s", strings.Join(errStr, ", "))
	}
	return rs, nil
}
