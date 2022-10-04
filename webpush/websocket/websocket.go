package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var ch chan int

func main() {
	var upgrader = &websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ch = make(chan int)
	g := gin.Default()
	g.GET("/set", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}, func(c *gin.Context) {
		time.Sleep(time.Second * 10)
		ch <- 0
	})
	g.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}, func(c *gin.Context) {
		w, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		<-ch
		if err != nil {
			panic(err)
		}
		go func() {
			w.WriteMessage(websocket.TextMessage, []byte("this is websocket"))
		}()
		_, p, err := w.ReadMessage()
		fmt.Println(string(p))
		defer w.Close()
	})
	g.Run(":8000")
}
