package util

import (
	"gopkg.in/ini.v1"
	"strconv"
	"sync"
)

type Config struct {
	lock     *sync.RWMutex
	cfg      *ini.File
	fileName string
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
func (cfg *Config) GetBooleanOrDefault(section, name string, defaultValue bool) bool {
	key, err := cfg.getSectionKey(section, name)
	if err != nil {
		return defaultValue
	} else {
		v := key.Value()
		if len(v) == 0 {
			return defaultValue
		}
		return EqualsAnyIgnoreCase(v, "true")
	}
}

func (cfg *Config) SetBoolean(section, key string, value bool) error {
	return cfg.SetString(section, key, BoolToString(value))
}
func (cfg *Config) SetString(section, key string, value string) error {
	sec := cfg.cfg.Section(section)
	if sec.HasKey(key) {
		preKey := sec.Key(key)
		preKey.SetValue(value)
		return nil
	} else {
		_, err := sec.NewKey(key, value)
		return err
	}
}
func (cfg *Config) SetInt(section, key string, value int) error {
	return cfg.SetString(section, key, strconv.Itoa(value))
}
func (cfg *Config) Save() error {
	return cfg.cfg.SaveTo(cfg.fileName)
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
	return &Config{lock: new(sync.RWMutex), cfg: cfg, fileName: fileName}, err
}
