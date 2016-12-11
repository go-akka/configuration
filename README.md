HOCON
=====

`example.go`

```go
package main

import (
	"fmt"
	"github.com/go-akka/configuration"
)

var configText = `
####################################
# Typesafe HOCON                   #
####################################

config {
  # Comment
  version = "0.0.1"
  one-second = 1s
  one-day = 1day
  array = ["one", "two", "three"] #comment
  bar = "bar"
  foo = foo.${config.bar} 
  number = 1
  object {
    a = "a"
    b = "b"
    c = {
        d = ${config.object.a} //comment
    }
}
`

func main() {
	conf := configuration.ParseString(configText)

	fmt.Println("config.one-second:", conf.GetTimeDuration("config.one-second"))
	fmt.Println("config.one-day:", conf.GetTimeDuration("config.one-day"))
	fmt.Println("config.array:", conf.GetStringList("config.array"))
	fmt.Println("config.bar:", conf.GetString("config.bar"))
	fmt.Println("config.foo:", conf.GetString("config.foo"))
	fmt.Println("config.number:", conf.GetInt64("config.number"))
	fmt.Println("config.object.a:", conf.GetString("config.object.a"))
	fmt.Println("config.object.c.d:", conf.GetString("config.object.c.d"))
}
```

```bash
> go run example.go
config.one-second: 1s
config.one-day: 24h0m0s
config.array: [one two three]
config.bar: bar
config.foo: foo.bar
config.number: 1
config.object.a: a
config.object.c.d: a
```
