package middlewares

import (
	apierrors "gitlab.cept.gov.in/it-2.0-common/n-api-errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		apierrors.HandleCommonError(c, err.Err)
	}
}
