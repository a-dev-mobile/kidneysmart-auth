package config

import (
	"errors"
	// Errors package for handling errors.
	"fmt"
	// Fmt package for formatting strings.
	"os"
	// Os package for interacting with the operating system, like file handling.
	"path/filepath"
	// Filepath package for manipulating filename paths.
	"gopkg.in/yaml.v3"
	// Yaml.v3 package for YAML processing.
)



// Config represents the entire configuration as structured in YAML.
type Config struct {
	Environment      Environment    `yaml:"environment"`
	Logging          LoggingConfig  `yaml:"logging"`
	ClientConnection ClientConfig   `yaml:"clientConnectionSettings"`
	ExternalService  ExternalConfig `yaml:"externalServiceIntegrations"`
	Database         DatabaseConfig `yaml:"database"`
}

type LoggingConfig struct {
	Level      string     `yaml:"level"`
	FileOutput FileConfig `yaml:"fileOutput"`
}

type FileConfig struct {
	FilePath       string         `yaml:"filePath"`
	RotationPolicy RotationPolicy `yaml:"rotationPolicy"`
	MaxSizeMB      int            `yaml:"maxSizeMB"`
	MaxBackups     int            `yaml:"maxBackups"`
}

type ClientConfig struct {
	GinMode        string   `yaml:"ginMode"`
	Port           string   `yaml:"port"`
	Host           string   `yaml:"host"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
}

type ExternalConfig struct {
	SmtpServer SmtpConfig `yaml:"smtpServer"`
}

type SmtpConfig struct {
	Grpc GrpcConfig `yaml:"grpc"`
}

type GrpcConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	User              string            `yaml:"user"`
	Password          string            `yaml:"password"`
	Host              string            `yaml:"host"`
	Port              string            `yaml:"port"`
	Name              string            `yaml:"name"`
	ConnectionTimeout int               `yaml:"connectionTimeoutSeconds"`
	MaxPoolSize       int               `yaml:"maxPoolSize"`
	Collections       map[string]string `yaml:"collections"`
}

// loadConfig reads and decodes the YAML configuration file.
// It is a private function, indicated by the lowercase first letter.
// Takes the file path as input and returns a pointer to the config struct or an error.
func LoadConfig(configPath string, configName string) (*Config, error) {
	configFile := filepath.Join(configPath, configName)

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		// Checks if the file exists. If it does not, returns an error.
		return nil, fmt.Errorf("config file does not exist: %s", configFile)
	}

	data, err := os.ReadFile(configFile)
	// Reads the file. If there is an error reading, it returns an error.
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Expand environment variables in the YAML
	expandedData := []byte(os.ExpandEnv(string(data)))
	// Declares a variable of type T to hold the configuration data.

	var config Config
	if err := yaml.Unmarshal(expandedData, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}
	// Returns a pointer to the config struct if successful.
	return &config, nil
}
