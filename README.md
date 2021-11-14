# GINITTER

##Just Rate Limiter
this middleware is for speed limiting in gin framework projects. Stay safe!

##Usage

### Rate Limiter

The request limit within *1 minute* is fixed as *1000*. A *status code of 500* is returned for requests after the *1000th* request.
```go
package main

import (
	"github.com/ahmetberke/ginitter"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	rl, err := ginitter.NewRateLimiter(ginitter.Config{
		Max: 1000,
		Expiration: time.Minute * 1,
		ExceededHandler: func(context *gin.Context) {
			context.AbortWithStatus(500)
			return
		},
	})
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(rl.Protect())
	r.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message":"hello gopherrs!!",
		})
	})
	err = r.Run()
}
```


### Rate Limiter With IP
The number of requests from an ip address in *1 minute* is limited to *10*. *A status code of 500* is returned for requests after the *10th* request

```go
package main

import (
	"github.com/ahmetberke/ginitter"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	rl, err := ginitter.NewRateLimiterIP(ginitter.Config{
		Max: 10,
		Expiration: time.Minute * 1,
		ExceededHandler: func(context *gin.Context) {
			context.AbortWithStatus(500)
			return
		},
	})
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(rl.Protect())
	r.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message":"hello gopherrs!!",
		})
	})
	err = r.Run()
}
```

