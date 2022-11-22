package configuration

import (
	"path/filepath"
	"strings"

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
	directory := filepath.Dir(path)
	viper.AddConfigPath(directory)
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	viper.SetConfigName(name)
	log.Infof("Attempting to load configuration from directory %s, name %s", directory, name)
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Failed to read in configuration: %s", err)
	}
	return err
}
