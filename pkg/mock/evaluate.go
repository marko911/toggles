package mock

import (
	"encoding/json"
	"io/ioutil"
	"toggle/server/pkg/errors"
	"toggle/server/pkg/models"
)

// Evaluate mocks the evaluate repository
type Evaluate struct {
	FlagPath string
}

// EvaluateInvalidFlagKey implements an evaluation repository with a receiver
// that is a mock
type EvaluateInvalidFlagKey struct{}

//GetFlag returns an error as a mock to a failed db request for a bad flag key
func (m *EvaluateInvalidFlagKey) GetFlag(key string) (*models.Flag, error) {
	return nil, errors.ErrFlagNotFound
}

//GetFlag returns a fixtured flag from JSON file
func (m *Evaluate) GetFlag(key string) (*models.Flag, error) {
	var flag models.Flag

	content, err := ioutil.ReadFile(m.FlagPath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(content, &flag)

	return &flag, nil
}

//EvaluateByte is a flag in JSON format
type EvaluateByte struct {
	Flag []byte
}

// GetFlag casts a flag from JSON
func (e *EvaluateByte) GetFlag(key string) (*models.Flag, error) {
	var flag models.Flag

	json.Unmarshal(e.Flag, &flag)
	return &flag, nil

}
