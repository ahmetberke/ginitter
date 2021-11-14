package ginitter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl, err := NewRateLimiter(Config{
		Max: 1,
		Expiration: time.Second * 3,
		ExceededHandler: func(context *gin.Context) {
			context.AbortWithStatus(500)
			return
		},
	})
	if err != nil {
		t.Errorf("error on creating new rate limiter")
	}

	r := gin.Default()
	r.Use(rl.Protect())
	r.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message":"hello gopherrs!!",
		})
	})
	err = r.Run()
	if err != nil {
		t.Errorf("err on running server")
	}
}

func TestNewRateLimiterIP(t *testing.T) {
	rl, err := NewRateLimiterIP(Config{
		Max: 1,
		Expiration: time.Second * 3,
		ExceededHandler: func(context *gin.Context) {
			context.AbortWithStatus(500)
			return
		},
	})
	if err != nil {
		t.Errorf("error on creating new rate limiter")
	}

	r := gin.Default()
	r.Use(rl.Protect())
	r.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message":"hello gopherrs!!",
		})
	})
	err = r.Run()
	if err != nil {
		t.Errorf("err on running server")
	}
}

func TestRateLimiter_Protect(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(20)
	counter := 0
	for i:=0;i<20;i++ {
		go func() {
			time.Sleep(2 * time.Second)
			resp, err := http.Get("http://localhost:8080")
			if err != nil {
				t.Errorf("Http get request error")
			}
			if resp.StatusCode == 200 {
				counter++
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if counter != 1 {
		t.Errorf("error on protect middleware, counter must be 1, counter: %v", counter)
	}
}