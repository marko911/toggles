package mock

import (
	"encoding/json"
	"io/ioutil"
	"toggle/server/pkg/models"

	"gopkg.in/mgo.v2/bson"
)

type mockRead struct {
	flagsJSON       []byte
	tenantJSON      []byte
	flagStatsJSON   []byte
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

	json.Unmarshal(m.tenantJSON, &t)

	return &t
}

func (m *mockRead) GetTenantFromAPIKey(apiKey string) *models.Tenant {
	var t models.Tenant

	json.Unmarshal(m.tenantJSON, &t)

	return &t
}

func (m mockRead) GetEvals() ([]models.Evaluation, error) {
	var e []models.Evaluation
	json.Unmarshal(m.evaluationsJSON, &e)
	return e, nil
}

func (m mockRead) GetFlagEvals(bson.ObjectId, int, int) ([]models.Evaluation, int, error) {
	var e []models.Evaluation
	json.Unmarshal(m.evaluationsJSON, &e)
	return e, 100, nil
}

func (m mockRead) GetFlagStats(id bson.ObjectId) (*models.FlagStats, error) {
	var s models.FlagStats
	json.Unmarshal(m.flagStatsJSON, &s)
	return &s, nil
}
