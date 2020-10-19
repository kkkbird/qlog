package qlog

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"math"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func getGinLogger(logger ...*logrus.Entry) *logrus.Entry {
	if len(logger) == 0 {
		return logrus.StandardLogger().WithField("pkg", "gin")
	}

	return logger[0]
}

// GinLogger is the qlog logger for GIN, copy from https://github.com/toorop/gin-logrus
func GinLogger(logger ...*logrus.Entry) gin.HandlerFunc {
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
			hostname = "unknown"
		}
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := log.WithFields(logrus.Fields{
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

type teeReadCloser struct {
	io.Reader
	io.Closer
}

func TeeReadCloser(rc io.ReadCloser, w io.Writer) io.ReadCloser {
	tee := io.TeeReader(rc, w)
	return &teeReadCloser{tee, rc}
}

type ginResponseMultiWriter struct {
	gin.ResponseWriter
	mw io.Writer
}

const (
	trimmedMsgLength = 16 * 1024
)

func GinMultiWriter(gw gin.ResponseWriter, w io.Writer) gin.ResponseWriter {
	mw := io.MultiWriter(gw, w)

	return &ginResponseMultiWriter{gw, mw}
}

func (w *ginResponseMultiWriter) Write(p []byte) (n int, err error) {
	return w.mw.Write(p)
}

func GinAPILogger(logger ...*logrus.Entry) gin.HandlerFunc {
	log := getGinLogger(logger...)

	return func(c *gin.Context) {
		start := time.Now()

		var bufReq, bufRsp bytes.Buffer

		c.Request.Body = TeeReadCloser(c.Request.Body, &bufReq)
		c.Writer = GinMultiWriter(c.Writer, &bufRsp)

		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		entry := log.WithFields(logrus.Fields{
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"userAgent":  c.Request.UserAgent(),
		})

		var msg bytes.Buffer
		msg.WriteString(c.Request.Method)
		msg.WriteString(" ")
		msg.WriteString(c.Request.RequestURI)

		reqLen := bufReq.Len()

		if reqLen > 0 {
			contentType := strings.ToLower(c.Request.Header.Get("Content-Type"))
			if strings.HasPrefix(contentType, "application/json") {
				msg.WriteString(" req:")

				if reqLen > trimmedMsgLength {
					msg.WriteString(fmt.Sprintf("[%d]", reqLen))
					msg.Write(bufReq.Bytes()[:trimmedMsgLength-12])
					msg.WriteString("...")
				} else {
					msg.Write(bufReq.Bytes())
				}
			} else {
				msg.WriteString(fmt.Sprintf(" req: len=%d, content_type=%s", reqLen, contentType))
			}
		}

		rspLen := bufRsp.Len()
		if rspLen > 0 {
			contentType := c.Writer.Header().Get("Content-Type")
			if strings.HasPrefix(contentType, "application/json") {
				msg.WriteString(" rsp:")

				if rspLen > trimmedMsgLength {
					msg.WriteString(fmt.Sprintf("[%d]", rspLen))
					msg.Write(bufRsp.Bytes()[:trimmedMsgLength-12])
					msg.WriteString("...")
				} else {
					msg.Write(bufRsp.Bytes())
				}
			} else {
				msg.WriteString(fmt.Sprintf(" rsp: len=%d, content_type=%s", rspLen, contentType))
			}
		}

		entry.Debug(msg.String())
	}
}
