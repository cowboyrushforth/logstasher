// Package logstasher is a Gin middleware that prints logstash-compatiable
// JSON to a given io.Writer for each HTTP request.
package logstasher

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strconv"
	"time"
)

type logstashEvent struct {
	Timestamp     string              `json:"@timestamp"`
	Version       int                 `json:"@version"`
	Method        string              `json:"method"`
	Path          string              `json:"path"`
	Status        int                 `json:"status"`
	Size          int                 `json:"size"`
	Duration      float64             `json:"duration"`
	Params        map[string][]string `json:"params,omitempty"`
	UserAgent     string              `json:"user_agent"`
	ContentLength string              `json:"content_length"`
	ContentType   string              `json:"content_type"`
	Referer       string              `json:"referer"`
	XForward      string              `json:"xforward"`
	Origin        string              `json:"origin"`
}

// Logger returns a middleware handler prints the request in a Logstash-JSON compatiable format
func Logger(writer io.Writer) gin.HandlerFunc {
	out := log.New(writer, "", 0)
	return func(c *gin.Context) {
		start := time.Now()

		rw := c.Writer
		c.Next()
		amount_written, _ := strconv.Atoi(rw.Header().Get("Content-Length"))
		params := make(map[string][]string)

		// GET params
		query_values := c.Request.URL.Query()
		for k, v := range query_values {
			params[k] = v
		}

		event := logstashEvent{
			time.Now().Format(time.RFC3339),
			1,
			c.Request.Method,
			c.Request.URL.Path,
			rw.Status(),
			amount_written,
			time.Since(start).Seconds() * 1000.0,
			params,
			c.Request.Header.Get("User-Agent"),
			c.Request.Header.Get("Content-Length"),
			c.Request.Header.Get("Content-Type"),
			c.Request.Header.Get("Referer"),
			c.Request.Header.Get("X-Forwarded-For"),
			c.Request.Header.Get("Origin"),
		}

		output, err := json.Marshal(event)
		if err != nil {
			// Should this be fatal?
			log.Printf("Unable to JSON-ify our event (%#v): %v", event, err)
			return
		}
		out.Println(string(output))
	}
}
