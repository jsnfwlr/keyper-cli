package app

import (
	env "github.com/caarlos0/env/v10"
	validator "github.com/go-playground/validator/v10"
)

// Config reads environment variables into the provided config struct.
// The struct field can be matched to the env var using the env tag.
// Defaults for the field can be set using the envDefault tag.
func Config(config any) error {
	err := env.Parse(config)
	if err != nil {
		return err
	}
	return nil
}

// Validate validates the provided config struct against the validator tags on the struct definition.
func Validate(config any) error {
	validate := validator.New()
	return validate.Struct(config)
}
