// Package sconfig is a utility to parse application configurations from
// environment variables and command line flags.
//
// This package uses https://github.com/spf13/viper under the hood and is
// compatible with CLI apps built with https://github.com/spf13/cobra.
package sconfig

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	spec       interface{}
	envEnabled bool
	envPrefix  string
	flagSet    *pflag.FlagSet
}

func New(spec interface{}) *Config {
	return &Config{spec: spec}
}

func (c *Config) FromEnvironment(envPrefix string) *Config {
	c.envEnabled = true
	c.envPrefix = envPrefix
	return c
}

func (c *Config) BindFlags(flagSet *pflag.FlagSet) *Config {
	c.flagSet = flagSet
	return c
}

func (c *Config) Parse() error {
	v := viper.New()

	if !isStructPointer(c.spec) {
		return ErrInvalidSpecification
	}

	err := c.setFields(v)
	if err != nil {
		return err
	}

	if c.envEnabled {
		// Automatically parse environment variables
		// Example: "MYAPP_PORT" will map to "Port",
		// where "MYAPP" is the envPrefix.
		v.SetEnvPrefix(c.envPrefix)
		v.AutomaticEnv()
	}

	err = v.Unmarshal(c.spec)
	if err != nil {
		return err
	}

	if c.envEnabled {
		err = c.checkRequiredFields()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) setFields(v *viper.Viper) error {
	return forEachStructField(c.spec,
		func(field reflect.StructField, value reflect.Value) error {
			def, ok := field.Tag.Lookup("default")
			if ok {
				v.SetDefault(field.Name, def)
			}

			if c.envEnabled {
				v.BindEnv(field.Name, "")
			}

			if c.flagSet != nil {
				if flag, ok := field.Tag.Lookup("flag"); ok {
					err := bindPFlag(
						v, c.flagSet, c.spec, flag, def, field, value)
					if err != nil {
						return &ErrInvalidField{
							Field: field.Name,
							Err:   err,
						}
					}
				}
			}

			return nil
		})
}

func (c *Config) checkRequiredFields() error {
	missingFields := make([]string, 0)

	forEachStructField(c.spec,
		func(field reflect.StructField, value reflect.Value) error {
			if required, ok := field.Tag.Lookup("required"); ok {
				if isTrue(required) && isZero(value) {
					missingFields = append(
						missingFields,
						envVarName(c.envPrefix, field.Name),
					)
				}
			}

			return nil
		})

	if len(missingFields) != 0 {
		return &ErrRequiredFields{
			Fields: missingFields,
		}
	}

	return nil
}

func isStructPointer(s interface{}) bool {
	p := reflect.ValueOf(s)
	if p.Kind() != reflect.Ptr {
		return false
	}

	v := p.Elem()
	return v.Kind() == reflect.Struct
}

func forEachStructField(
	s interface{},
	f func(reflect.StructField, reflect.Value) error,
) error {
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		err := f(field, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func envVarName(prefix, name string) string {
	if prefix != "" {
		return strings.ToUpper(prefix + "_" + name)
	}

	return strings.ToUpper(name)
}

func isZero(value reflect.Value) bool {
	return reflect.DeepEqual(
		value.Interface(),
		reflect.Zero(reflect.TypeOf(value.Interface())).Interface(),
	)
}
