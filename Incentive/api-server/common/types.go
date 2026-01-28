package common

import (
	"gitlab.cept.gov.in/it-2.0-common/n-api-server/util/wrapper"

	"github.com/gin-gonic/gin"
)

type (
	MiddlewareGroup = []gin.HandlerFunc
	GinAppWrapper   = wrapper.Wrapper[*gin.Engine]
)
