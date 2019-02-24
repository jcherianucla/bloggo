package config

import (
	"io/ioutil"

	"go.uber.org/fx"
	"gopkg.in/yaml.v2"
)

const (
	_base = "base.yaml"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In
}

type Result struct {
	fx.Out

	AppConfig
}

type AppConfig interface {
	Config() *Config
}

type Config struct {
	Service    string     `yaml:"service"`
	GRPCConfig GRPCConfig `yaml:"grpc"`
}

func (c *Config) Config() *Config {
	return c
}

type GRPCConfig struct {
	HostPort string `yaml:"hostPort"`
}

func New(p Params) (Result, error) {
	cfg := &Config{}
	yamlFile, err := ioutil.ReadFile(_base)
	if err != nil {
		return Result{}, err
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return Result{}, err
	}
	return Result{AppConfig: cfg}, nil
}
