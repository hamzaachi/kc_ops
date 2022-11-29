package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Kc_source kc_instance `yaml:"kc_source"`
	Kc_target kc_instance `yaml:"kc_target"`
}

type kc_instance struct {
	Url      string   `yaml:"url"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Realm    string   `yaml:"realm"`
	Clients  []string `yaml:"clients",omitempty`
	Roles    []string `yaml:"roles",omitempty`
	Groups   []string `yaml:"groups",omitempty`
}

func New(path string) (s *Config, err error) {
	s = &Config{}
	yamlReader, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Error reading config file: %s", err)
		return nil, err
	}
	defer yamlReader.Close()

	decoder := yaml.NewDecoder(yamlReader)
	decoder.KnownFields(true)
	if err = decoder.Decode(s); err != nil {
		fmt.Errorf("error parsing config file: %s", err)
		return nil, err
	}

	return s, nil
}

// TODO
func validator(c *Config) {
	fmt.Println(c)
}
