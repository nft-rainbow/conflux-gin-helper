package ginutils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func DataResponse(data interface{}) interface{} {
	return data
}

func ErrorResponse(code int, err error) gin.H {
	return gin.H{
		"code":    code,
		"message": err.Error(),
	}
}

func SuccessResponse() gin.H {
	return gin.H{
		"code":    0,
		"message": "Success",
	}
}

type List[T any] struct {
	Count int64 `json:"count"`
	Items []T   `json:"items"`
}

func RenderResp(c *gin.Context, data interface{}, err error) {
	if err != nil {
		RenderRespError(c, err, 500)
		return
	}
	RenderRespOK(c, data)
}

func RenderListResp[T any](c *gin.Context, count int64, items []T, err error) {
	if err != nil {
		RenderRespError(c, err, 500)
		return
	}
	RenderRespOK(c, List[T]{Count: count, Items: items})
}

func RenderRespOK(c *gin.Context, data interface{}, httpStatusCode ...int) {
	statusCode := http.StatusOK

	if len(httpStatusCode) > 0 {
		statusCode = httpStatusCode[0]
	}
	c.JSON(statusCode, DataResponse(data))
}

// rainbowErrorCode 有值时，message 为 err message，如果 err 为 rainbow error, 则 status code 与 code 都来自 err, 否则来自rainbowErrorCode
// 否则 message 为 err message，status code 与 code 为 ERR_INTERNAL_SERVER_COMMON
func RenderRespError(c *gin.Context, err error, status int, code ...int) {
	c.Error(err)
	c.Set("error_stack", fmt.Sprintf("%+v", errors.WithStack(err)))

	_rainbowErrorCode := status * 100

	if len(code) > 0 {
		_rainbowErrorCode = code[0]
	}

	c.JSON(status, ErrorResponse(int(_rainbowErrorCode), err))
}

type textErr string

func (e textErr) Error() string {
	return string(e)
}

func RenderRespErrorText(c *gin.Context, err string, status int, code ...int) {
	RenderRespError(c, textErr(err), status, code...)
}
