package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var lockKey = "mylock"        // 你的锁的名称
var lockValue = "mylockValue" // 你的锁的名称
var redisClient = NewFoxRedis("127.0.0.1:6379", "")
var mutex = &sync.Mutex{}
var wg sync.WaitGroup
var successCount = 0
var failureCount = 0

func TestFoxRedis_Acquire(t *testing.T) {
	numTests := 50 // 设置要运行的并发测试数量

	for i := 0; i < numTests; i++ {
		wg.Add(1)
		go testLock()
	}

	wg.Wait()
	fmt.Printf("Successful locks: %d\n", successCount)
	fmt.Printf("Failed locks: %d\n", failureCount)
	return
}

func testLock() {
	defer wg.Done()
	if redisClient.AcquireLockWithRetry(lockKey, lockValue, 3*time.Second, 3, 1*time.Second) {
		successCount++
		// 执行一些需要锁的操作
		time.Sleep(800 * time.Millisecond)
		redisClient.ReleaseLock(lockKey, lockValue)
	} else {
		failureCount++
	}
}
