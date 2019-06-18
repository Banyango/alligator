package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	CacheSize int
	Proxy []Proxy
}

type Rule struct {
	Type    string
	Pattern []string
}

func (rule *Rule) String() string {
	return fmt.Sprint("Matcher { Type = ", rule.Type, " } Patterns => ", rule.Pattern)
}

type Proxy struct {
	Name  string
	Path  string
	Host  string
	Scheme string
	Rules []Rule
}

func New(fileContents string) (*Config, error) {
	var config Config

	if _, err := toml.Decode(fileContents, &config); err != nil {
		return nil, err
	}

	config.setDefaults()

	return &config, nil
}

func (config *Config) setDefaults() {
	if config.CacheSize == 0 {
		log.Println("Defaulting CacheSize to 10MB")
		config.CacheSize = 10*1024*1024
	}
}
