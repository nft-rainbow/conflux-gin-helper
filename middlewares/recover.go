package middlewares

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/conflux-gin-helper/utils/ginutils"
	"github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	var buf bytes.Buffer
	return gin.CustomRecoveryWithWriter(&buf, gin.RecoveryFunc(func(c *gin.Context, err interface{}) {
		defer func() {
			fmt.Println(buf.String())
			logrus.WithField("recovered", buf.String()).WithField("error", err).Error("panic and recovery")
			buf.Reset()
		}()
		ginutils.RenderRespError(c, errors.New("server error"), 500)
		c.Abort()
	}))
}
