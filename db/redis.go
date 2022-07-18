package db

import (
	"context"
	redis "github.com/go-redis/redis/v9"
	"fmt"
)

var Rdb *redis.Client
var ctx = context.TODO()

func init() {
	Rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	fmt.Println(Rdb.Ping(ctx).Result())
}