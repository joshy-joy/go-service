package pdf

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Address string `yaml:"address"`
}

type APIEntry struct {
	Method string `yaml:"method"`
	URL    string `yaml:"url"`
}

type Config struct {
	Server ServerConfig        `yaml:"server"`
	API    map[string]APIEntry `yaml:"-"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		svcPath := filepath.Join("resource", "config", "service.yaml")
		apiPath := filepath.Join("resource", "config", "api.yaml")
		sb, err := os.ReadFile(svcPath)
		if err != nil {
			panic(err)
		}
		var c Config
		if err := yaml.Unmarshal(sb, &c); err != nil {
			panic(err)
		}
		ab, err := os.ReadFile(apiPath)
		if err != nil {
			panic(err)
		}
		var api map[string]APIEntry
		if err := yaml.Unmarshal(ab, &api); err != nil {
			panic(err)
		}
		c.API = api
		cfg = &c
	})
	return cfg
}
