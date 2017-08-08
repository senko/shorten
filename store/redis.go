package store

import (
	"gopkg.in/redis.v3"
	"log"
	"strconv"
	"sync"
)

type Redis struct {
	redis     *redis.Client
	keyPrefix string
	lock      sync.Mutex
}

type RedisOptions struct {
	RedisAddr string
	KeyPrefix string
}

const (
	defaultRedisAddr = ":6379"
	defaultKeyPrefix = "shorten"
	urlMapSuffix     = "-urls"
	reverseMapSuffix = "-reverse-urls"
	hitMapSuffix     = "-hits"
)

func NewRedis(options *RedisOptions) *Redis {
	redisAddr := options.RedisAddr
	if redisAddr == "" {
		redisAddr = defaultRedisAddr
	}
	keyPrefix := options.KeyPrefix
	if keyPrefix == "" {
		keyPrefix = defaultKeyPrefix
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	return &Redis{
		redis:     client,
		keyPrefix: keyPrefix,
	}
}

func (s *Redis) GetFullUrl(shortUrl string) string {
	res, err := s.redis.HGet(s.keyPrefix+reverseMapSuffix, shortUrl).Result()
	if err == redis.Nil {
		return "" // key doesn't exist
	}
	if err != nil {
		log.Printf("store.GetFullUrl(%s): error getting full url from redis: %s", shortUrl, err)
		return ""
	}
	return res
}

func (s *Redis) GetShortUrl(fullUrl string) string {
	res, err := s.redis.HGet(s.keyPrefix+urlMapSuffix, fullUrl).Result()
	if err == redis.Nil {
		return "" // key doesn't exist
	}
	if err != nil {
		log.Printf("store.GetShortUrl(%s): error getting short url from redis: %s", fullUrl, err)
		return ""
	}
	return res
}

func (s *Redis) SetUrl(url, shortUrl string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	ok, err := s.redis.HSetNX(s.keyPrefix+reverseMapSuffix, shortUrl, url).Result()
	if err != nil {
		log.Printf("store.Redis.SetUrl(%s): error saving short->full mapping to redis: %s", url, err)
		return false
	}
	if !ok {
		return false
	}

	_, err = s.redis.HSet(s.keyPrefix+urlMapSuffix, url, shortUrl).Result()
	if err != nil {
		log.Printf("store.Redis.SetUrl(%s): error saving full->short mapping to redis: %s", url, err)
		s.redis.HDel(s.keyPrefix+reverseMapSuffix, shortUrl).Result()
		return false
	}
	return true
}

func (s *Redis) CountUrls() int {
	res, err := s.redis.HLen(s.keyPrefix + urlMapSuffix).Result()
	if err != nil {
		log.Printf("store.Redis.CountUrls(): error counting URLs in redis: %s", err)
		return 0
	}
	return int(res)
}

func (s *Redis) FlushUrls() {
	_, err := s.redis.Del(s.keyPrefix + urlMapSuffix).Result()
	if err != nil {
		log.Printf("store.Redis.FlushUrls(): error removing all full->short mappings from redis: %s", err)
	}
	_, err = s.redis.Del(s.keyPrefix + reverseMapSuffix).Result()
	if err != nil {
		log.Printf("store.Redis.FlushUrls(): error removing all short->full mappings from redis: %s", err)
	}
}

func (s *Redis) RecordHit(url string) {
	_, err := s.redis.HIncrBy(s.keyPrefix+hitMapSuffix, url, 1).Result()
	if err != nil {
		log.Printf("store.Redis.SetUrl(%s): error recording hit to redis: %s", url, err)
	}
}

func (s *Redis) GetUrlHits(url string) int {
	res, err := s.redis.HGet(s.keyPrefix+hitMapSuffix, url).Result()
	if err != nil {
		log.Printf("store.Redis.GetUrlHits(%s): error getting hits from redis: %s", url, err)
	}
	n, err := strconv.Atoi(res)
	if err != nil {
		log.Printf("store.Redis.GetUrlHits(%s): invalid value '%s': %s", url, res, err)
		return 0
	}
	return n
}

func (s *Redis) GetAllHits() map[string]int {
	res, err := s.redis.HGetAllMap(s.keyPrefix + hitMapSuffix).Result()
	if err != nil {
		log.Printf("store.Redis.GetAllHits(): error getting hits from redis: %s", err)
		return map[string]int{}
	}
	hits := make(map[string]int, len(res))
	for k, v := range res {
		if n, err := strconv.Atoi(v); err == nil {
			hits[k] = n
		}
	}
	return hits
}

func (s *Redis) FlushHits() {
	_, err := s.redis.Del(s.keyPrefix + hitMapSuffix).Result()
	if err != nil {
		log.Printf("store.Redis.FlushUrls(): error removing all URLs from redis: %s", err)
	}
}
