package helpers

import (
	"ScrabShortener/cache"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = cache.RedisConnection()

func GetViewCount() any {
	iter := rdb.Scan(ctx, 0, "*", 0).Iterator()

	newViews := 1000

	for iter.Next(ctx) {
		key := iter.Val()

		viewCount, err := rdb.HGet(ctx, key, "view_count").Int()
		if err != nil && err != redis.Nil {
			log.Fatal(err, "1")
			continue
		}

		prevViewCount, err := rdb.HGet(ctx, key, "prev_view_count").Int()
		if err != nil && err != redis.Nil {
			log.Fatal(err, "2")
			continue
		}

		if viewCount >= prevViewCount+newViews {
			err := rdb.Expire(ctx, key, 10*time.Minute).Err()
			if err != nil && err != redis.Nil {
				log.Fatal(err, "3")
				continue
			}
		}

		err = rdb.HSet(ctx, key, "prev_view_count", viewCount).Err()
		if err != nil && err != redis.Nil {
			log.Fatal(err, "4")
			continue
		}
	}
	return nil
}

func TickerFunc() {

	ticker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			GetViewCount()
		}
	}
}
