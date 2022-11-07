package cache

import (
	"challenge/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis/v9"
)

type Cache struct {
	client *redis.Client
}

const REDIS_ADDR_ENV = "REDIS_PORT"

//abrir conexao com o banco do redis
func NewClient() (*Cache, error) {
	cl, ok := os.LookupEnv(REDIS_ADDR_ENV)
	if !ok {
		return nil, fmt.Errorf("fail to get redis address")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cl,
		Password: "",
		DB: 0,
	})
	return &Cache{client: redisClient}, nil
}

func (c *Cache) Get( key string) ( *entity.Task, error ){
	var tk *entity.Task
	tc, err := c.client.Get(context.Background(), key).Result()
	if  err == nil{
		json.Unmarshal([]byte(tc), tk)
		return tk, nil
	}

	return tk, err    

}

func (c *Cache) Set( key string, value string) error { 
	return c.client.Set(context.Background(), key, value, 0).Err()
}

func (c *Cache) Del( key string) error {
	return c.client.Del(context.Background(), key).Err()
}


