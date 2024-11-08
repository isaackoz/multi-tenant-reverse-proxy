package mycache

import (
	"net/http"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
)

type RateLimitCache = *cache.Cache[int]

func CreateRateLimitCacheManager() RateLimitCache {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)
	return cache.New[int](ristrettoStore)
}

func RateLimitByUserIp(r *http.Request) {

}
