package mock

import (
	"encoding/json"
	"io/ioutil"
	"toggle/server/pkg/models"
)

type mockRead struct {
	flagsJSON       []byte
	tenantJSON      []byte
	singleFlagJSON  []byte
	evaluationsJSON []byte
	segsPath        *string
	usersPath       *string
}

func (m *mockRead) GetFlags(models.Tenant) ([]models.Flag, error) {
	var flags []models.Flag

	json.Unmarshal(m.flagsJSON, &flags)

	return flags, nil
}

func (m *mockRead) GetFlag(key string) (*models.Flag, error) {
	var flag models.Flag
	json.Unmarshal(m.singleFlagJSON, &flag)
	return &flag, nil
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

func (m *mockRead) GetTenant(key string) *models.Tenant {
	var t models.Tenant

	json.Unmarshal(m.flagsJSON, &t)

	return &t
}

func (m mockRead) GetEvals() ([]models.Evaluation, error) {
	var e []models.Evaluation
	json.Unmarshal(m.evaluationsJSON, &e)
	return e, nil
}
