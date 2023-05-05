package cache

import (
	"encoding/json"
	"time"

	"git.01.alem.school/bbaktyke/test.project.git/internal/models"

	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) *redisCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value *models.ProblemWithTopics) error {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	client.Set(key, json, cache.expires*time.Minute)
	return nil
}

func (cache *redisCache) Get(key string) *models.ProblemWithTopics {
	client := cache.getClient()
	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}
	problem := models.ProblemWithTopics{}
	err = json.Unmarshal([]byte(val), &problem)
	if err != nil {
		return nil
	}

	return &problem
}

func (cache *redisCache) SetArr(key string, value *[]models.ProblemWithTopics) error {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	client.Set(key, json, cache.expires*time.Minute)
	return nil
}

func (cache *redisCache) GetArr(key string) *[]models.ProblemWithTopics {
	client := cache.getClient()
	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}
	problem := []models.ProblemWithTopics{}
	err = json.Unmarshal([]byte(val), &problem)
	if err != nil {
		return nil
	}

	return &problem
}
