package configuration

import (
	"github.com/snasphysicist/ferp/v2/pkg/configuration/router"
	"github.com/snasphysicist/ferp/v2/pkg/log"
)

// populateMethodRouters adds method routers for all configured redirects and forwarded routes
func populateMethodRouters(c Configuration) (Configuration, error) {
	ii, iiErr := populateMethodRoutersForIncomings(c.HTTP.Incoming)
	c.HTTP.Incoming = ii
	ird, irdErr := populateMethodRoutersForRedirects(c.HTTP.Redirects)
	c.HTTP.Redirects = ird
	// TODO: HTTPS
	err := joinNonNilErrors([]error{iiErr, irdErr}, ", ", "invalid methods: %s")
	return c, err
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
	err := joinNonNilErrors(errs, ", ", "invalid methods: %s")
	return iswr, err
}

// populateMethodRoutersForRedirects populates the method routers fields in the provided
// redirects, returning an error summarising which (if any) had invalid methods
func populateMethodRoutersForRedirects(rds []Redirect) ([]Redirect, error) {
	rdswr := make([]Redirect, 0)
	errs := make([]error, 0)
	for _, rd := range rds {
		mrs, err := findMethodRouters(rd.Methods)
		errs = append(errs, err)
		rdswr = append(rdswr, Redirect{
			From:          rd.From,
			To:            rd.To,
			Methods:       rd.Methods,
			MethodRouters: mrs,
		})
		log.Infof("For %+v, method routers %+v", rd, mrs)
	}
	err := joinNonNilErrors(errs, ", ", "invalid methods: %s")
	return rdswr, err
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
	return rs, joinNonNilErrors(errs, ", ", "%s")
}
