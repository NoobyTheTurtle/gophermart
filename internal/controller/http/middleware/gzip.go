package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func shouldCompress(r *http.Request, contentType string) bool {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}

	compressibleTypes := []string{
		"text/",
		"application/json",
		"application/javascript",
	}

	for _, t := range compressibleTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}

	return false
}

func GzipMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.String(http.StatusBadRequest, "Invalid gzip body")
				c.Abort()
				return
			}
			c.Request.Body = reader
			defer reader.Close()
		}

		contentType := c.Writer.Header().Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
		}

		if !shouldCompress(c.Request, contentType) {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			logger.Error("Failed to create gzip writer", "error", err)
			c.String(http.StatusInternalServerError, "internal server error")
			c.Abort()
			return
		}
		defer gz.Close()

		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Del("Content-Length")

		c.Writer = gzipWriter{ResponseWriter: c.Writer, Writer: gz}
		c.Next()
	}
}
