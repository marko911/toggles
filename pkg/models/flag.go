package models

import "toggle/server/pkg/errors"

// Validate validates a flag payload from POST request
func (f *Flag) Validate() (bool, error) {
	if f.Key == "" || f.Name == "" || len(f.Variations) == 0 {
		return false, errors.ErrJSONPayloadInvalidFlag
	}
	for _, variation := range f.Variations {
		if variation.Name == "" {
			return false, errors.ErrJSONPayloadInvalidFlag
		}
	}
	return true, nil

}
