package config

import (
	"io/ioutil"

	"go.uber.org/fx"
	"gopkg.in/yaml.v2"
)

const (
	_base = "base.yaml"
)

// Module provides the Config through Fx
var Module = fx.Provide(New)

// Params defines the input dependencies for the Config
type Params struct {
	fx.In
}

// Result defines the output dependency that is the Application config
type Result struct {
	fx.Out
	AppConfig
}

// AppConfig defines an interface to interact with the service configuration
// as created through the yaml
type AppConfig interface {
	// Config provides singleton access to the configuration
	Config() *Config
}

// Config houses the underlying configuration as separate objects
type Config struct {
	Service       string        `yaml:"service"`
	GRPCConfig    GRPCConfig    `yaml:"grpc"`
	LoggerConfig  LoggerConfig  `yaml:"logger"`
	MetricsConfig MetricsConfig `yaml:"metrics"`
	DBConfig      DBConfig      `yaml:"datastore"`
}

func (c *Config) Config() *Config {
	return c
}

// GRPCConfig holds the configuration for the grpc server
type GRPCConfig struct {
	// HostPort is the ip_addr:port the server should run on
	HostPort string `yaml:"hostPort"`
}

// LoggerConfig holds the configuration for the logger
type LoggerConfig struct {
}

// MetricsConfig hold the configuration for the metrics client
type MetricsConfig struct {
	// Port is the port that the prometheus client should ping for metrics
	Port string `yaml:"port"`
}

// DBConfig holds the configuration for the underlying data store
type DBConfig struct {
	// User is the db user that performs the operation
	User string `yaml:"user"`
	// Password is the password for the user's access to the db
	Password string `yaml:"password"`
	// Name is the name for the data base
	Name string `yaml:"name"`
	// HostPort is the ip_addr:port the data store should run on
	HostPort string `yaml:"hostPort"`
}

// New returns a new application configuration converting the yaml to
// corresponding configuration objects through their tags
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
