package mock

import (
	"context"
	"testing"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	"toggle/server/pkg/read"
	"toggle/server/pkg/store/mongo"
)

// CreateContext takes context wrapping functions and returns a context with all functions applied to it
func CreateContext(t *testing.T, f ...func(c context.Context) context.Context) context.Context {
	ctx := context.TODO()

	ctx = initCreateService(ctx)
	ctx = initReadService(ctx)
	ctx = initDatabaseMemoryStore(ctx)
	ctx = initCache(ctx)

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
		flagsJSON: []byte(`[
			{
				"_id": { "$oid": "5f09cd037815899375759b9b" },
				"name": "Early testers",
				"key": "",
				"enabled": true,
				"variations": [
					{ "name": "On", "percent": 100 },
					{ "name": "Off", "percent": 0 }
				],
				"tenant": { "$oid": "5ef5f06a4fc7eb0006772c49" }
			},
			{
				"_id": { "$oid": "5f09d08d40a5b800068a5d88" },
				"name": "Young chicks",
				"key": "hey-ladies",
				"enabled": true,
				"variations": [
					{ "name": "On", "percent": 100 },
					{ "name": "Off", "percent": 0 }
				],
				"targets": [
					{
						"rules": [
							{ "attribute": "gender", "operator": "EQ", "value": "female" }
						],
						"variations": [
							{ "name": "On", "percent": 100 },
							{ "name": "Off", "percent": 0 }
						]
					}
				],
				"tenant": { "$oid": "5ef5f06a4fc7eb0006772c49" }
			}
		]
		`),
	})
	return context.WithValue(c, read.ServiceKey, mockReadService)
}

func initCreateService(c context.Context) context.Context {
	mockCreateService := create.NewService(&mockCreate{})
	return context.WithValue(c, create.ServiceKey, mockCreateService)
}

func initCache(c context.Context) context.Context {
	cache := auth.GetCache()
	return context.WithValue(c, auth.CacheServiceKey, cache)
}
