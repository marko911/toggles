package mock

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

type MockManager struct {
	controller *gomock.Controller
}

func CreateContext(t *testing.T, f ...func(c context.Context) context.Context) (context.Context, *MockManager) {
	ctx := context.TODO()
	ctx = initDatabaseMemoryStore(ctx)

}

func initDatabaseMemoryStore(c context.Context) context.Context {
	mockStore := mockManager.Store()
	return context.WithValue(c, store.CTXKey, mockStore)
}
