package cache

import (
	"github.com/dgraph-io/ristretto"
	jsoniter "github.com/json-iterator/go"
)

type Cache struct {
	*ristretto.Cache
}

func New(maxKeys, maxCost int64) (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxKeys,
		MaxCost:     maxCost,
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
