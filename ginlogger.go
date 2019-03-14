package qlog

import (
	"fmt"
	"time"

	"math"
	"os"

	"github.com/gin-gonic/gin"
)

func getGinLogger(logger ...*Entry) *Entry {
	if len(logger) == 0 {
		return StandardLogger().WithField("pkg", "gin")
	}

	return logger[0]
}

// GinLogger is the qlog logger for GIN, copy from https://github.com/toorop/gin-logrus
func GinLogger(logger ...*Entry) gin.HandlerFunc {
	log := getGinLogger(logger...)

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := log.WithFields(Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%d \"%s %s\" (%dms)", statusCode, c.Request.Method, path, latency)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Trace(msg)
			}
		}
	}
}
