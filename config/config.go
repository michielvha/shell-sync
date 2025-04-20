package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SyncIntervalSec int             `yaml:"sync_interval_sec"`
	Backend         BackendConfig   `yaml:"backend"`
	HistoryFiles    []HistoryConfig `yaml:"history_files"`
	Filter          FilterConfig    `yaml:"filter"`
	LogFile         string          `yaml:"log_file"`
}

type BackendConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type HistoryConfig struct {
	Path string `yaml:"path"`
	Shell string `yaml:"shell"`
}

type FilterConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Patterns []string `yaml:"patterns"`
	Action   string   `yaml:"action"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.SyncIntervalSec == 0 {
		cfg.SyncIntervalSec = 15
	}
	return &cfg, nil
}
