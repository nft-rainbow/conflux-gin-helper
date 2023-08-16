package middlewares

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/conflux-gin-helper/utils/ginutils"
)

func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")
		sizeStr := c.DefaultQuery("limit", "10")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			ginutils.RenderRespError(c, errors.New("invalid pagenation"), 400)
		}
		limit, err := strconv.Atoi(sizeStr)
		if err != nil {
			ginutils.RenderRespError(c, errors.New("invalid pagenation"), 400)
		}
		// apply default value
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		c.Set("page", page)
		c.Set("limit", limit)
		c.Set("offset", (page-1)*limit)
		c.Next()
	}
}
