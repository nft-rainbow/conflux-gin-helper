package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogOptions struct {
	IgnoredPaths     []string
	ReqHeaderLogger  func(headers http.Header) interface{}
	RespHeaderLogger func(headers http.Header) interface{}
}

func Logger(logOption *LogOptions) gin.HandlerFunc {
	if logOption == nil {
		logOption = &LogOptions{}
	}

	var _ignoredPaths = make(map[string]bool)
	for _, v := range logOption.IgnoredPaths {
		_ignoredPaths[v] = true
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// record body
		var body []byte
		if !_ignoredPaths[strings.ToLower(path)] && c.Request.ContentLength < 5*1024 {
			body, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		}

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		entry := logrus.WithFields(logrus.Fields{
			"status code":    param.StatusCode,
			"latency":        fmt.Sprintf("%13v", param.Latency),
			"client ip":      fmt.Sprintf("%15s", param.ClientIP),
			"method":         param.Method,
			"path":           param.Path,
			"full path":      c.FullPath(),
			"body":           string(body),
			"content length": c.Request.ContentLength,
		})

		if logOption.ReqHeaderLogger != nil {
			entry = entry.WithField("req header", logOption.ReqHeaderLogger(c.Request.Header))
		}

		if logOption.RespHeaderLogger != nil {
			entry = entry.WithField("resp header", logOption.ReqHeaderLogger(c.Writer.Header()))
		}

		if param.ErrorMessage != "" {
			entry = entry.
				WithField("errors", param.ErrorMessage).
				WithField("stack", c.GetString("error_stack"))

			if c.GetString("error_stack") != "" {
				fmt.Printf("Request error %v\n%v\n", param.ErrorMessage, c.GetString("error_stack"))
			}
		}
		entry.Info("Request")
	}
}
