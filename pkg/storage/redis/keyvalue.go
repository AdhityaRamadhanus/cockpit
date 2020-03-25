package redis

import (
	"time"

	"github.com/AdhityaRamadhanus/cockpit"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

//KeyValueService implements cockpit.KeyValueService interface using redis
type KeyValueService struct {
	redisClient *redis.Client
}

//NewKeyValueService construct a new KeyValueService from redis client
func NewKeyValueService(redisClient *redis.Client) *KeyValueService {
	return &KeyValueService{
		redisClient: redisClient,
	}
}

//Get a cache in bytes from a key
func (c KeyValueService) Get(key string) (result []byte, err error) {
	val, err := c.redisClient.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, cockpit.ErrKeyNotFound
		}
		return nil, errors.Wrapf(err, "redisClient.Get(%q).Bytes() err", key)
	}

	return val, nil
}

//Set cache in bytes with key without expiration
func (c KeyValueService) Set(key string, value []byte) (err error) {
	if err := c.redisClient.Set(key, string(value), 0).Err(); err != nil {
		return errors.Wrapf(err, "redisClient.Set(%q, <val>, 0) err", key)
	}

	return nil
}

func (c KeyValueService) GetHashAll(key string) (map[string]string, error) {
	hash, err := c.redisClient.HGetAll(key).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, cockpit.ErrKeyNotFound
		}
		return nil, errors.Wrapf(err, "redisClient.HGetAll(%q).Result() err", key)
	}

	return hash, nil
}

func (c KeyValueService) SetHashAll(key string, value map[string]interface{}) error {
	if err := c.redisClient.HMSet(key, value).Err(); err != nil {
		return errors.Wrapf(err, "redisClient.HMSet(%q, value) err", key)
	}

	return nil
}

//Delete cache in bytes with key without expiration
func (c KeyValueService) Delete(key string) (err error) {
	if err := c.redisClient.Del(key).Err(); err != nil {
		return errors.Wrapf(err, "redisClient.Delete(%q) err", key)
	}

	return nil
}

//SetEx cache in bytes with key with expiration
func (c KeyValueService) SetEx(key string, value []byte, expiration time.Duration) (err error) {
	if err := c.redisClient.Set(key, value, expiration).Err(); err != nil {
		return errors.Wrapf(err, "redisClient.Set(%q, <val>, %d) err", key, expiration)
	}

	return nil
}
