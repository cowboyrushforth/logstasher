# logstasher 

logstasher is a gin middleware that prints logstash-compatible JSON to an `io.Writer` for each HTTP request.

Here's an example from one of the Go microservices we have at @bikeexchange :

``` json
{
  "@timestamp":"2014-03-01T19:08:06+11:00","@version":1,"method":"GET",
  "path":"/locations/slugs/VIC/Williams-Landing","status":200,"size":238,
  "duration":14.059902000000001,"params":{"country":["au"]}
}
```

Used in conjunction with the [rotating file writer](http://github.com/mipearson/rfw) it allows for rotatable logs ready to feed directly into logstash with no parsing.

### Example

``` go
package main

import (
  "log"

  "github.com/gin-gonic/gin"
  "github.com/mipearson/logstasher"
  "github.com/mipearson/rfw"
)

func main() {
  r := gin.Default()

  logstashLogFile, err := rfw.Open("hello.log", 0644)
  if err != nil {
    log.Fatalln(err)
  }
  defer logstashLogFile.Close()
  r.Use(logstasher.Logger(logstashLogFile))

  r.Get("/", func(c *gin.Context) string {
    c.String(200, "hello world")
  })
  r.Run()
}
```

```
## logstash.conf
input {
  file {
    path => ["hello.log"]
    codec => "json"
  }
}
```
