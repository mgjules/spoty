package cache

import (
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/json"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/fx"
)

const _defaultBufferItems = 64

// Module exported for initialising a new Cache.
var Module = fx.Options(
	fx.Provide(New),
)

// Cache is a simple wrapper around ristretto.Cache.
type Cache struct {
	*ristretto.Cache
}

// New creates a new Cache.
func New(cfg *config.Config) (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.CacheMaxKeys,
		MaxCost:     cfg.CacheMaxCost,
		BufferItems: _defaultBufferItems,
		Cost: func(value interface{}) int64 {
			test, err := json.Marshal(value)
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
