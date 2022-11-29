package configuration

import (
	"fmt"

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
	err := joinNonNilErrors(errs, ", ", "invalid mapper configuration: %s")
	return c, err
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
	jErr := joinNonNilErrors(errs, ", ",
		fmt.Sprintf("no mapper matching configuration %#v (failed to match: %s)", c, "%s"))
	return nil, jErr
}
