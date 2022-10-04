package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var ch chan int

func main() {
	g := gin.Default()
	ch = make(chan int)
	flag := true
	g.GET("/set", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}, func(c *gin.Context) {
		go func() {
			time.Sleep(10 * time.Second)
			ch <- 0
		}()
		go func() {
			<-ch
			flag = false
		}()
	})
	g.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}, func(c *gin.Context) {
		if !flag {
			c.String(http.StatusOK, "this is short poll")
		}
		c.String(http.StatusOK, "timeout")
	})
	g.Run(":8000")
}
