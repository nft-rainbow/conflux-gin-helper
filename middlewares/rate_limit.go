package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var RateLimitMiddleware gin.HandlerFunc

func InitRateLimitMiddleware(limitPerSec int64) {
	var rate = limiter.Rate{
		Period: 1 * time.Second,
		Limit:  limitPerSec,
	}
	var store = memory.NewStore()
	var instance = limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	RateLimitMiddleware = mgin.NewMiddleware(instance)

	logrus.WithField("limit config", rate).Info("set limit config")
}
