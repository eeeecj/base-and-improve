package main

import (
	"fmt"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var ch chan int

func main_() {
	g := gin.Default()
	ch = make(chan int)
	g.GET("/set", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}, func(c *gin.Context) {
		go func() {
			time.Sleep(3 * time.Second)
			ch <- 0
		}()
	})
	g.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
	}, func(c *gin.Context) {
		fmt.Println("waiting")
		<-ch
		fmt.Println("hello")
		f, ok := c.Writer.(http.Flusher)
		if !ok {
			panic("writer is not a flusher")
		}
		fmt.Fprintf(c.Writer, "data:this is sse\n\n")
		f.Flush()
		fmt.Println("close")
	})
	g.Run(":8000")
}

func main() {
	g := gin.Default()
	ch = make(chan int)
	g.GET("/set", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")

	}, func(c *gin.Context) {
		go func() {
			time.Sleep(3 * time.Second)
			ch <- 0
		}()
	})
	g.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
	}, func(c *gin.Context) {
		<-ch
		sse.Encode(c.Writer, sse.Event{
			Data: "this is sse",
		})
		fmt.Println("close")
	})
	g.Run(":8000")
}
