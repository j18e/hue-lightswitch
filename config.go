package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ghodss/yaml"
)

type Config struct {
	Switches []*Switch  `json:"switches"`
	Mappings []*Mapping `json:"mappings"`
}

type Switch struct {
	Name    string    `json:"name"`
	Model   string    `json:"model"`
	Buttons []*Button `json:"buttons"`
}

type Button struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Mapping struct {
	Switch string `json:"switch"`
	Button string `json:"button"`
	Lights []int  `json:"lights"`

	// Payload is the JSON payload to be sent to Hue bridge.
	Payload string `json:"payload"`

	LastFired time.Time `json:"-"`
}

func NewConfig(file string) (*Config, error) {
	bs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(bs, &cfg); err != nil {
		return nil, err
	}

	if len(cfg.Switches) == 0 {
		return nil, errors.New("no switches configured")
	}
	if len(cfg.Mappings) == 0 {
		return nil, errors.New("no mappings configured")
	}
	for _, m := range cfg.Mappings {
		if !cfg.btnInCfg(m.Switch, m.Button) {
			return nil, fmt.Errorf("parsing mappings: button %s/%s not found", m.Switch, m.Button)
		}
	}

	return &cfg, nil
}

func (cfg *Config) btnInCfg(sw, btn string) bool {
	for _, s := range cfg.Switches {
		if s.Name == sw {
			for _, b := range s.Buttons {
				if b.Name == btn {
					return true
				}
			}
		}
	}
	return false
}
