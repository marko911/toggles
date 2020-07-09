package mock

import "toggle/server/pkg/models"

type mockCreate struct {
}

func (m *mockCreate) InsertFlag(*models.Flag) error {
	return nil
}
func (m *mockCreate) InsertSegment(*models.Segment) error {
	return nil
}
func (m *mockCreate) InsertUser(*models.User) error {
	return nil
}
func (m *mockCreate) InsertAttributes([]models.Attribute) error {

	return nil
}