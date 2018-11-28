package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	_, err := redis.String(redisConn.Do("SET", "key", "lock", "PX", 11, "NX"))
	if err != nil {
		fmt.Println(err)
	}
}
