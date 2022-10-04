package scheduler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Run() {
	app := gin.Default()
	app = Route(app)
	if err := app.Run(":9909"); err != nil {
		fmt.Println(err)
	}
}
