package middlewares

import (
	"net/http"
	//r "testencrypt/rate-gin/ratelimiter"

	rate "gitlab.cept.gov.in/it-2.0-common/n-api-server/ratelimiter"

	"github.com/gin-gonic/gin"
)

func RateMiddleware(globalBucket *rate.LeakyBucket) gin.HandlerFunc {
	return func(c *gin.Context) {

		if globalBucket.Allow() {
			c.Next()
		} else {

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Traffic shaping limit exceeded",
			})
		}

	}
}
