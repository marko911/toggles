package mock

import "toggle/server/pkg/models"

type mockCreate struct {
}

func (m *mockCreate) CreateFlag(*models.Flag) error {
	return nil
}
func (m *mockCreate) CreateSegment(*models.Segment) error {
	return nil
}
func (m *mockCreate) CreateUser(*models.User) error {
	return nil
}
func (m *mockCreate) CreateAttributes(*models.User) error {
	return nil
}
