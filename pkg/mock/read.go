package mock

import (
	"encoding/json"
	"io/ioutil"
	"toggle/server/pkg/models"
)

type mockRead struct {
	flagsJSON []byte
	segsPath  *string
	usersPath *string
}

func (m *mockRead) GetFlags(models.Tenant) ([]models.Flag, error) {
	var flags []models.Flag

	json.Unmarshal(m.flagsJSON, &flags)

	return flags, nil
}

func (m *mockRead) GetSegments(models.Tenant) ([]models.Segment, error) {
	var segments []models.Segment
	content, err := ioutil.ReadFile(*m.segsPath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(content, &segments)

	return segments, nil
}

func (m *mockRead) GetUsers(models.Tenant) ([]models.User, error) {
	var users []models.User
	content, err := ioutil.ReadFile(*m.usersPath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(content, &users)

	return users, nil
}
