package scheduler

import "github.com/gin-gonic/gin"

func Route(app *gin.Engine) *gin.Engine {
	jobController := NewJobController()
	app.GET("/jobs", jobController.GetJobs)
	app.POST("/jobs", jobController.AddJob)
	app.DELETE("/jobs/:id", jobController.DeleteJob)
	return app
}
