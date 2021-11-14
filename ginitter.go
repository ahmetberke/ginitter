package ginitter

import (
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

type Config struct {
	Max 			int
	Expiration 		time.Duration
	ExceededHandler gin.HandlerFunc
}

type rateLimiter struct {
	max int
	expiration 		time.Duration
	currentCount 	int
	exceededHandler gin.HandlerFunc
	mutex 			*sync.Mutex
}

var defaultExceededHandler gin.HandlerFunc = func(context *gin.Context) {
	context.AbortWithStatus(500)
	return
}

func NewRateLimiter(config Config) (*rateLimiter, error) {

	mutex := sync.Mutex{}
	if config.ExceededHandler == nil {
		config.ExceededHandler = defaultExceededHandler
	}
	rl := &rateLimiter{
		max: config.Max,
		expiration: config.Expiration,
		exceededHandler: config.ExceededHandler,
		mutex: &mutex,
	}

	go func() {
		for true {
			time.Sleep(rl.expiration)
			rl.currentCount = 0
		}
	}()

	return  rl, nil
}

func (rl *rateLimiter) Protect() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := rl.add()
		if err != nil {
			rl.exceededHandler(context)
		}
	}
}

func (rl *rateLimiter) add() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	if rl.currentCount >= rl.max {
		return errors.New("limit expected")
	}
	rl.currentCount++
	return nil
}

type rateLimiterIP struct {
	max 			int
	expiration 		time.Duration
	currentIPs		map[string]int
	exceededHandler gin.HandlerFunc
	mutex 			*sync.Mutex
}

func NewRateLimiterIP(config Config) (*rateLimiterIP, error) {
	currentIPs := make(map[string]int)
	mutex := sync.Mutex{}
	if config.ExceededHandler == nil {
		config.ExceededHandler = defaultExceededHandler
	}
	rl := &rateLimiterIP{
		max: config.Max,
		expiration: config.Expiration,
		exceededHandler: config.ExceededHandler,
		mutex: &mutex,
		currentIPs: currentIPs,
	}

	go func() {
		for true {
			time.Sleep(rl.expiration)
			rl.currentIPs = make(map[string]int)
		}
	}()

	return  rl, nil
}

func (rl *rateLimiterIP) Protect() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := rl.add(context.ClientIP())
		if err != nil {
			rl.exceededHandler(context)
		}
	}
}

func (rl *rateLimiterIP) add(ip string) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	if rl.currentIPs[ip] >= rl.max {
		return errors.New("limit expected")
	}
	rl.currentIPs[ip]++
	return nil
}