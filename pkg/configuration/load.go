package configuration

import (
	"github.com/mitchellh/mapstructure"
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
	return c, nil
}

// readInConfiguration configures viper to target the configuration file and attempts to read it in
func readInConfiguration(path string) error {
	// TODO: implement
	return nil
}
