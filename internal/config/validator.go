package config

import "github.com/pkg/errors"

func validateAppEnvironment(v string) (interface{}, error) {
	switch v {
	case "local":
		return Local, nil
	case "development":
		return Development, nil
	case "production":
		return Production, nil
	default:
		return nil, errors.Errorf("unknown environment: %s", v)
	}
}
