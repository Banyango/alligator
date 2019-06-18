package config

import "github.com/BurntSushi/toml"

type Config struct {
	Proxy []Proxy
}

type Rule struct {
	Type    string
	Pattern []string
}

type Proxy struct {
	Name  string
	Path  string
	Host  string
	Scheme string
	Rules []Rule
}

func NewConfig(fileContents string) (*Config, error) {
	var config Config

	if _, err := toml.Decode(fileContents, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
