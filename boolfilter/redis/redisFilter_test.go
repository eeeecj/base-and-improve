package redis

import (
	"context"
	"github.com/balance/boolfilter"
	"github.com/go-redis/redis/v8"
	"strconv"
	"testing"
)

func TestCacheFilter(t *testing.T) {
	var options = &redis.Options{
		Addr:     "127.0.0.1:6379",
		Username: "",
		Password: "",
		DB:       0,
	}
	var key = "test"
	cli := redis.NewClient(options)
	cachedFilter, err := NewCacheFilter(context.TODO(), "boolm", key, cli, 10024, boolfilter.DefaultHash...)
	if err != nil {
		t.Fatal(err)
	}
	for i := 200; i <= 300; i++ {
		cachedFilter.Push([]byte(strconv.Itoa(i)))
	}
	t.Log(cachedFilter.Exists([]byte(strconv.Itoa(290)))) // true
	t.Log(cachedFilter.Exists([]byte(strconv.Itoa(299)))) // true
	t.Log(cachedFilter.Exists([]byte(strconv.Itoa(350)))) // false
	// must use write to save the data to redis.
	// 必须使用write方法将数据全部提交到redis服务器
	cachedFilter.Write()
	cachedFilter, err = NewCacheFilter(context.TODO(), "boolm", key, cli, 10024, boolfilter.DefaultHash...)
	t.Log(cachedFilter.Exists([]byte(strconv.Itoa(290)))) // true
	t.Log(cachedFilter.Exists([]byte(strconv.Itoa(299)))) // true
	t.Log(cachedFilter.Exists([]byte(strconv.Itoa(350)))) // false
}
