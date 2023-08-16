package middlewares

import (
	"io"
	"log"
	runtimePporf "runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/conflux-gin-helper/utils"
)

func SlowMonitor(slowWriter io.Writer, traceWriter io.Writer) gin.HandlerFunc {
	slowLog := log.New(slowWriter, "[Slow Request]\t", 0)
	return func(c *gin.Context) {
		done := false
		go func() {
			requestTime := time.Now()
			<-time.After(10 * time.Second)
			if !done {
				utils.DingWarn("find very slow request, %v %v", c.Request.Method, c.Request.URL.Path)
				slowLog.Printf("%v, %v, %v, %v\n", requestTime.Format(time.StampMilli), c.ClientIP(), c.Request.Method, c.Request.URL.Path)
				runtimePporf.Lookup("mutex").WriteTo(slowWriter, 1)
				logTrace(traceWriter)
			}
		}()
		c.Next()
		done = true
	}
}

func logTrace(traceWriter io.Writer) {
	// traceFilePath := path.Join(config.GetConfig().Log.Folder, fmt.Sprintf("trace_%d.out", time.Now().UnixMilli()))
	// f, err := os.OpenFile(traceFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	// if err != nil {
	// 	panic(err)
	// }
	trace.Start(traceWriter)
	defer trace.Stop()
	<-time.After(10 * time.Second)
}
