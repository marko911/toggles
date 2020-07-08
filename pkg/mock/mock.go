package mock

import (
	"context"
	"testing"
	"toggle/server/pkg/create"
	"toggle/server/pkg/read"
	"toggle/server/pkg/store/mongo"

	"github.com/golang/mock/gomock"
)

type MockManager struct {
	controller *gomock.Controller
}

func CreateContext(t *testing.T, f ...func(c context.Context) context.Context) (context.Context, *MockManager) {
	ctx := context.TODO()
	ctx = initDatabaseMemoryStore(ctx)
	ctx = initReadService(ctx)
}

func initDatabaseMemoryStore(c context.Context) context.Context {
	mockStore := mongo.NewMockStore()
	return context.WithValue(c, mongo.CTXKey, mockStore)
}

func initReadService(c context.Context) context.Context {
	mockReadService := read.NewService(&mockRead{}) // TODO:add paths, figure out how to pass in path
	return context.WithValue(c, read.ServiceKey, mockReadService)
}

func initCreateService(c context.Context) context.Context {
	mockCreateService := create.NewService(&mockCreate{})
}
