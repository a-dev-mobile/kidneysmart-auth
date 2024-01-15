package config

import (
	"fmt"
	"strings"
)

// UnmarshalYAML custom unmarshalling for Environment.
func (e *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var envStr string
	if err := unmarshal(&envStr); err != nil {
		return err
	}

	switch Environment(envStr) {
	case Dev, Prod:
		*e = Environment(envStr)
		return nil
	default:
		return fmt.Errorf("invalid environment: %s", envStr)
	}
}

// UnmarshalYAML custom unmarshalling for RotationPolicy.
func (r *RotationPolicy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var policyStr string
	if err := unmarshal(&policyStr); err != nil {
		return err
	}

	switch RotationPolicy(policyStr) {
	case Monthly, Weekly, Daily:
		*r = RotationPolicy(policyStr)
		return nil
	default:
		return fmt.Errorf("invalid rotation policy: %s", policyStr)
	}
}

// UnmarshalYAML customizes the unmarshalling for LogLevel.
func (l *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var levelStr string
	if err := unmarshal(&levelStr); err != nil {
		return err
	}

	levelStr = strings.ToLower(levelStr)
	switch LogLevel(levelStr) {
	case LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError:
		*l = LogLevel(levelStr)
		return nil
	default:
		return fmt.Errorf("invalid log level: %s", levelStr)
	}
}