package logstasher

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
        "github.com/gin-gonic/gin"
        "fmt"
)

// Copied and adapted from martini's logger_test.go
func Test_Logger(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

        r := gin.Default()
        r.Use(Logger(buff))
	r.Use(func(c *gin.Context) {
             c.String(404, "not found")
             c.Next()
	})

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar?baz=quux", nil)
	//req.ParseForm()
	if err != nil {
		t.Error(err)
	}

        r.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)

        fmt.Println("BUFF", buff.String())
	refute(t, len(buff.String()), 0)

	var event logstashEvent
	err = json.Unmarshal(buff.Bytes(), &event)
	if err != nil {
		t.Error(err)
	}
	expect(t, event.Version, 1)
	expect(t, event.Method, "GET")
	expect(t, event.Path, "/foobar")
	expect(t, event.Status, http.StatusNotFound)
	refute(t, event.Duration, 0)
	refute(t, event.Timestamp, "")
	expect(t, event.Size, 0)
	expect(t, event.Params, map[string][]string{"baz": []string{"quux"}})
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Got %#v, expected %#v", line, a, b)
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Got %#v, which was not expected", line, a)
	}
}
