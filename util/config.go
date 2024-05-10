package util

import (
	"gopkg.in/ini.v1"
	"sync"
)

type Config struct {
	lock *sync.RWMutex
	cfg  *ini.File
}

func (cfg *Config) GetString(section, name string) string {
	key, err := cfg.getSectionKey(section, name)
	if err != nil {
		return ""
	} else {
		return key.Value()
	}
}
func (cfg *Config) GetStringOrDefault(section, name string, defaultValue string) string {
	key, err := cfg.getSectionKey(section, name)
	if err != nil {
		return defaultValue
	} else {
		v := key.Value()
		if len(v) == 0 {
			return defaultValue
		}
		return v
	}
}
func (cfg *Config) GetInt(section, name string) (int, error) {
	key, err := cfg.getSectionKey(section, name)
	if err != nil {
		return 0, err
	} else {
		return key.Int()
	}
}
func (cfg *Config) GetIntOrDefault(section, name string, defaultValue int) int {
	key, err := cfg.getSectionKey(section, name)
	if err != nil {
		return defaultValue
	} else {
		v, err := key.Int()
		if err != nil {
			return defaultValue
		}
		return v
	}
}

func (cfg *Config) GetInt64OrDefault(section, name string, defaultValue int64) int64 {
	key, err := cfg.getSectionKey(section, name)
	if err != nil {
		return defaultValue
	} else {
		v, err := key.Int64()
		if err != nil {
			return defaultValue
		}
		return v
	}
}

func (cfg *Config) getSectionKey(section, name string) (*ini.Key, error) {
	sc, err := cfg.cfg.GetSection(section)
	if err != nil {
		return nil, err
	} else {
		return sc.GetKey(name)
	}
}

func LoadFile(fileName string) (*Config, error) {
	cfg, err := ini.Load(fileName)
	if err != nil {
		return nil, err
	}
	return &Config{lock: new(sync.RWMutex), cfg: cfg}, err
}
