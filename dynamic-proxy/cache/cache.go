package mycache

import (
	"context"
	"dynamic-proxy/config"
	"dynamic-proxy/db"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
	"github.com/redis/go-redis/v9"
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

type TenantCache = *cache.LoadableCache[string]

func CreateTenantCacheManager(config *config.Config) TenantCache {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})

	fmt.Println("Creating tenant cache manager\n")

	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	// Init stores
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)
	redisStore := redis_store.NewRedis(redisClient, store.WithExpiration(5*time.Second))

	// Fallback load function for when data is not in either cache
	// It will fetch it from Postgres
	loadFunction := func(ctx context.Context, key any) (string, error) {
		strKey, ok := key.(string)
		if !ok {
			return "", fmt.Errorf("key is not a string")
		}
		tenant, err := db.GetTenant(strKey)

		if err != nil {
			return "", err
		}

		return tenant.ID, nil
	}

	cacheChain := cache.NewChain[string](
		cache.New[string](ristrettoStore),
		cache.New[string](redisStore),
	)

	loadableStore := cache.NewLoadable[string](loadFunction, cacheChain)

	/*
		First the reistrettoStore will be checked (which is local Go memory)
		If that fails, then the redisStore will be checked (which is the Redis database, which can be either local or external)
		If those both fail, we will run the loadFunction which is going to get the tenant from the PostgreSQL database

		The cache chain will automatically send data back up the chain to set the previously failed caches.
		For example, if both ristretto and redis miss, but the loadFunction is succesful, the cache chain will send the value from the
		load function and set it in both ristretto and redis so the next query will be a hit.

		If all miss, an err is returned

	*/

	return loadableStore
}
