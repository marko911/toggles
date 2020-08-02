package models

import "toggle/server/pkg/errors"

// Validate validates a flag payload from POST request
func (f *Flag) Validate() (bool, error) {
	if f.Key == "" {
		return false, errors.ErrJSONPayloadInvalidFlagKey
	}

	if f.Name == "" {
		return false, errors.ErrJSONPayloadInvalidName
	}
	if len(f.Variations) == 0 {
		return false, errors.ErrJSONPayloadInvalidVariations
	}
	for _, variation := range f.Variations {
		if variation.Name == "" {
			return false, errors.ErrJSONPayloadInvalidVariationEmpty
		}
	}
	return true, nil

}

// HasLimit tells us if flag has a limit of number of evaluations
func (f *Flag) HasLimit() bool {
	return f.Limit > 0
}
