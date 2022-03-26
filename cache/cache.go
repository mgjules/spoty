package cache

import (
	"github.com/JulesMike/spoty/config"
	"github.com/dgraph-io/ristretto"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(New),
)

type Cache struct {
	*ristretto.Cache
}

func New(cfg *config.Config) (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.CacheMaxKeys,
		MaxCost:     cfg.CacheMaxCost,
		BufferItems: 64,
		Cost: func(value interface{}) int64 {
			test, err := jsoniter.Marshal(value)
			if err != nil {
				return 1
			}
			return int64(len(test))
		},
	})
	if err != nil {
		return nil, err
	}

	return &Cache{cache}, nil
}
