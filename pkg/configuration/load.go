package configuration

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/snasphysicist/ferp/v2/pkg/functional"
	"github.com/snasphysicist/ferp/v2/pkg/log"
	"github.com/spf13/viper"
)

// Load attempts to load and validate the server configuration from the given path
func Load(path string) (Configuration, error) {
	err := readInConfiguration(path)
	if err != nil {
		return Configuration{}, err
	}
	c := Configuration{}
	err = viper.Unmarshal(&c, func(c *mapstructure.DecoderConfig) { c.TagName = "config" })
	if err != nil {
		log.Errorf("Failed to deserialise configuration: %s", err)
		return Configuration{}, err
	}
	log.Infof("Loaded and deserialised configuration: %#v", c)
	c, err = validate(c)
	if err != nil {
		log.Errorf("The configuration is not valid: %s", err)
		return Configuration{}, err
	}
	return c, nil
}

// readInConfiguration configures viper to target the configuration file and attempts to read it in
func readInConfiguration(path string) error {
	directory := filepath.Dir(path)
	viper.AddConfigPath(directory)
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	viper.SetConfigName(name)
	log.Infof("Attempting to load configuration from directory %s, name %s", directory, name)
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Failed to read in configuration: %s", err)
	}
	return nil
}

// validate ensures that all options provided in the configuration are valid
func validate(c Configuration) (Configuration, error) {
	c, pmErr := populatePathMappers(c)
	c, dErr := populateDownstreams(c)
	c, mrErr := populateMethodRouters(c)
	errs := functional.Map(
		functional.Filter([]error{pmErr, dErr, mrErr},
			func(e error) bool { return e != nil }),
		func(e error) string { return e.Error() })
	if len(errs) > 0 {
		return c, fmt.Errorf("invalid configuration: %s", strings.Join(errs, ", "))

	}
	return c, nil
}
