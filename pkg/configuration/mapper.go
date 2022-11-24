package configuration

import (
	"fmt"
	"strings"

	"github.com/snasphysicist/ferp/v2/pkg/functional"
	"github.com/snasphysicist/ferp/v2/pkg/mapper"
)

// populatePathMappers attempts to find and instantiate a path mapper
// for all downstream path mapper configurations
func populatePathMappers(c Configuration) (Configuration, error) {
	ds := make([]Downstream, 0)
	errs := make([]error, 0)
	for _, d := range c.Downstreams {
		m, err := loadPathMapper(d.MapperData)
		errs = append(errs, err)
		d.Mapper = m
		ds = append(ds, d)
	}
	c.Downstreams = ds
	errs = functional.Filter(errs, func(e error) bool { return e != nil })
	errStr := functional.Map(errs, func(e error) string { return e.Error() })
	if len(errs) > 0 {
		return c, fmt.Errorf(
			"invalid mapper configuration: %s", strings.Join(errStr, ", "))
	}
	return c, nil
}

// loadPathMapper attempts to find and instantiate a path mapper
// matching the provided configuration map
// TODO: all extensibility with custom path mappers added at "runtime"
func loadPathMapper(c map[string]string) (pathMapper, error) {
	mappers := []pathMapper{&mapper.Passthrough{}, &mapper.RemovePrefix{}}
	errs := make([]error, 0)
	for _, m := range mappers {
		err := m.From(c)
		if err == nil {
			return m, nil
		}
		errs = append(errs, err)
	}
	errStr := functional.Map(errs, func(e error) string { return e.Error() })
	return nil, fmt.Errorf(
		"no mapper matching configuration %#v (failed to match: %s)",
		c, strings.Join(errStr, ", "))
}
