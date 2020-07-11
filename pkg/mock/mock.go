package mock

import (
	"context"
	"testing"
	"toggle/server/pkg/create"
	"toggle/server/pkg/read"
	"toggle/server/pkg/store/mongo"

	"github.com/golang/mock/gomock"
)

// Manager is a mock manager instance
type Manager struct {
	controller *gomock.Controller
}

// Mock returns a new gomock controller
func Mock(t *testing.T) *Manager {
	mockCtrl := gomock.NewController(t)
	return &Manager{
		controller: mockCtrl,
	}
}

//Done signals to controller that mocks are called
func (m *Manager) Done() {
	m.controller.Finish()
}

// CreateContext takes context wrapping functions and returns a context with all functions applied to it
func CreateContext(t *testing.T, f ...func(c context.Context) context.Context) context.Context {
	ctx := context.TODO()

	ctx = initCreateService(ctx)
	ctx = initReadService(ctx)
	ctx = initDatabaseMemoryStore(ctx)

	for idx := range f {
		ctx = f[idx](ctx)

	}
	return ctx

}

func initDatabaseMemoryStore(c context.Context) context.Context {
	mockStore := mongo.NewMockStore()
	return context.WithValue(c, mongo.CTXKey, mockStore)
}

func initReadService(c context.Context) context.Context {
	mockReadService := read.NewService(&mockRead{
		flagsPath: "../../config/flags.json",
	})
	return context.WithValue(c, read.ServiceKey, mockReadService)
}

func initCreateService(c context.Context) context.Context {
	mockCreateService := create.NewService(&mockCreate{})
	return context.WithValue(c, create.ServiceKey, mockCreateService)
}
