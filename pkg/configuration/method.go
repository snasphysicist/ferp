package configuration

import (
	"fmt"
	"strings"

	"github.com/snasphysicist/ferp/v2/pkg/configuration/router"
	"github.com/snasphysicist/ferp/v2/pkg/functional"
)

// validateAllMethods ensure that all methods across the configuration are valid/supported
func validateAllMethods(c Configuration) error {
	httpErrs := functional.Map(
		functional.Filter(
			functional.Map(c.HTTP.Redirects,
				func(r Redirect) error { return validateMethods(r.Methods, r) }),
			func(e error) bool { return e != nil }),
		func(e error) string { return e.Error() })
	if len(httpErrs) != 0 {
		return fmt.Errorf("invalid methods in HTTP configuration: %s", strings.Join(httpErrs, ", "))
	}
	return nil
}

// validateMethods validates that the configured HTTP methods in the given incoming are valid options
func validateMethods(ms []string, source interface{}) error {
	errs := functional.Map(
		functional.Filter(
			functional.Map(ms, func(m string) error { return validateMethod(m) }),
			func(e error) bool { return e != nil }),
		func(e error) string { return e.Error() })
	if len(errs) != 0 {
		return fmt.Errorf("invalid methods in %#v: %s", source, strings.Join(errs, ", "))
	}
	return nil
}

// validateMethod validates that the provided HTTP methods is a valid option
func validateMethod(m string) error {
	if !functional.Contains(validMethods(), m) {
		return fmt.Errorf("invalid method %s", m)
	}
	return nil
}

// validMethods lists all valid HTTP method options in the incoming configurations.
// The special value "ALL" means proxy requests for all method types.
func validMethods() []string {
	ms := make([]string, 0)
	for m := range router.MethodRouters() {
		ms = append(ms, m)
	}
	return ms
}
