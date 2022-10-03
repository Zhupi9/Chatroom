package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func initPool(addr string, maxIdle, maxAct int, timeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxAct,
		IdleTimeout: timeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
	fmt.Println("Redis Pool connected: ", pool)
}
