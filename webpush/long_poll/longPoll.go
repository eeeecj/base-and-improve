package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var ch chan int

func main() {
	g := gin.Default()
	ch = make(chan int)
	g.GET("/set", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}, func(c *gin.Context) {
		go func() {
			time.Sleep(10 * time.Second)
			ch <- 0
		}()
	})
	g.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}, func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		select {
		case <-ch:
			c.String(http.StatusOK, "this is long poll")
		case <-ctx.Done():
			c.String(http.StatusOK, "timeout")
		}
	})
	g.Run(":8000")
}
