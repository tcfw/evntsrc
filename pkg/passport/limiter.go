package passport

import (
	"net"
	"os"
	"time"

	"github.com/go-redis/redis"
)

const (
	userMaxCount = 10
	userLimitTTL = time.Minute * 5
	ipLimitTTL   = time.Minute * 60
	ipMaxCount   = 100
)

func clearRateLimit(email string, ip net.IP) error {
	client, err := redisClient()
	if err != nil {
		return err
	}

	userKey := userCountKey(email, ip)
	ipKey := ipCountKey(ip)

	count, err := client.Del(userKey, ipKey).Result()
	if err != nil || count != 2 {
		return err
	}

	return nil
}

type keyTTL struct {
	key string
	ttl time.Duration
}

func incRateLimit(email string, ip net.IP) {
	keysToIncrement := []keyTTL{
		keyTTL{key: userCountKey(email, ip), ttl: userLimitTTL},
		keyTTL{key: ipCountKey(ip), ttl: ipLimitTTL},
	}

	client, err := redisClient()
	if err != nil {
		return
	}
	defer client.Close()

	for _, limit := range keysToIncrement {
		exists := len(client.Keys(limit.key).Val())
		if exists == 0 {
			client.Set(limit.key, 0, limit.ttl)
		}
		client.Incr(limit.key)
	}
}

func checkIPRateLimit(ip net.IP) (bool, time.Duration, int) {
	key := ipCountKey(ip)

	return checkKeyRateLimit(key, ipMaxCount)
}

func checkUserRateLimit(username string, ip net.IP) (bool, time.Duration, int) {
	key := userCountKey(username, ip)

	return checkKeyRateLimit(key, userMaxCount)
}

func checkKeyRateLimit(key string, limit int) (bool, time.Duration, int) {
	client, err := redisClient()
	if err != nil {
		return true, 0, limit //Allow on error
	}
	defer client.Close()

	exists := len(client.Keys(key).Val())

	if exists == 0 {
		return true, 0, limit //Allow if count doesn't exist
	}

	count, err := client.Get(key).Int64()
	if err != nil {
		return true, 0, limit //Allow on error
	}

	ttl := client.TTL(key).Val()
	if int(count) <= limit {
		return true, ttl, limit - int(count)
	}

	return false, ttl, limit
}

func userCountKey(email string, ip net.IP) string {
	return "login_rate:" + ip.String() + ":" + email
}

func ipCountKey(ip net.IP) string {
	return "login_rate:" + ip.String()
}

func redisClient() (*redis.Client, error) {
	redisConnHost, exists := os.LookupEnv("REDIS_HOST")
	if exists == false {
		redisConnHost = "redis:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisConnHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
