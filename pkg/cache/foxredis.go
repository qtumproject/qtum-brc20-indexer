package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type FoxRedis struct {
	Client *redis.Client
}

var ctx = context.Background()

func NewFoxRedis(addr, pass string) *FoxRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		panic("connected to redis failed [ERROR]=> " + err.Error())
	}
	return &FoxRedis{
		Client: client,
	}
}

// ReleaseLock
// key: lock key
// value: lock value, like tx_id, to make sure the lock to be release is lock by you
func (redis FoxRedis) ReleaseLock(key string, value string) {
	// 使用 Lua 脚本释放锁，确保原子性
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	_, err := redis.Client.Eval(ctx, luaScript, []string{key}, value).Result()
	if err != nil {
		log.Printf("Error releasing lock: %v", err)
	}
}

func (redis FoxRedis) AcquireLock(key, value string) bool {
	return redis.AcquireLockWithRetry(key, value, 3*time.Second, 0, 0)
}

// AcquireLockWithRetry acquire lock, if fail to acquire lock, it will retry maxRetries times after sleep(retryInterval)
func (redis FoxRedis) AcquireLockWithRetry(key string, value string, expireTime time.Duration, maxRetries int, retryInterval time.Duration) bool {
	attempt := 0
	for attempt <= maxRetries {
		// 尝试获取锁
		lockAcquired := redis.AcquireLockWithExpireTime(key, value, expireTime)
		if lockAcquired {
			//fmt.Printf("get lock success in attempt %d\n", attempt)
			return true
		}
		//fmt.Printf("get lock failed in attempt %d\n", attempt)
		attempt++
		if attempt <= maxRetries && retryInterval != 0 {
			time.Sleep(retryInterval)
		}
	}
	//fmt.Printf("get lock failed\n")
	return false
}

// AcquireLockWithExpireTime get lock with specified expiration time
func (redis FoxRedis) AcquireLockWithExpireTime(key, value string, expireTime time.Duration) bool {
	// 使用 SETNX 命令尝试获取锁
	lockAcquired, err := redis.Client.SetNX(ctx, key, value, expireTime).Result()
	if err != nil {
		log.Fatalf("Error acquiring lock: %v", err)
	}
	return lockAcquired
}

// SetMarshalValue marshal the struct data into []byte and then store
func (redis FoxRedis) SetMarshalValue(ctx context.Context, key string, value interface{}, expireTime time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redis.Client.Set(ctx, key, jsonValue, expireTime).Err()
}

func (redis FoxRedis) GetUnmarshalValue(ctx context.Context, key string, value interface{}) error {
	storedData, err := redis.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(storedData), &value)
}
