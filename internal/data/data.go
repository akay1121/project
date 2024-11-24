// Package data mainly focuses on the implementation of the repository interface defined in the biz package.
// It encapsulates the interaction with db and cache where it reads the stored data and maps it to the
// domain object defined in the biz package.
package data

import (
	"github.com/redis/go-redis/v9"
	"project/internal/conf"
	"project/internal/ent"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewCache,
	NewProjectRepository,
)

// Data wraps the db client
type Data struct {
	Client *ent.Client
}

// Cache wraps the Redis client
type Cache struct {
	Client *redis.Client
}

// NewData establishes the connection to the db based on the configuration
func NewData(c *conf.Data) (data *Data, cleanup func(), err error) {
	var dbClient *ent.Client
	if dbClient, err = ent.Open(c.Database.Driver, c.Database.Source); err != nil {
		log.Error(err)
	}
	cleanup = func() {
		log.Info("closing the data resources")
		if err := dbClient.Close(); err != nil {
			log.Error(err)
		}
	}
	data = &Data{Client: dbClient}
	return
}

// NewCache establishes the connection to the Redis server based on the configuration
func NewCache(c *conf.Data) (cache *Cache, cleanup func(), err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Network:      c.Redis.Network,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})
	cleanup = func() {
		log.Info("closing the connection to Redis server")
		if err := rdb.Close(); err != nil {
			log.Error(err)
		}
	}
	cache = &Cache{Client: rdb}
	return
}
