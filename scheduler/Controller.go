package scheduler

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"net/http"
	"strconv"
)

type JobControllerImp struct {
}

func (j *JobControllerImp) GetJobs(c *gin.Context) {
	var results []map[string]interface{}
	for _, e := range Cron.Entries() {
		results = append(results, map[string]interface{}{
			"id":   e.ID,
			"next": e.Next,
		})
	}
	c.JSON(http.StatusOK, results)
}

func (j *JobControllerImp) AddJob(c *gin.Context) {
	var payload struct {
		Cron string `json:"cron"`
		Exec string `json:"exec"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := Cron.AddFunc(payload.Cron, func() {
		ExecuteTask(payload.Exec)
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (j *JobControllerImp) DeleteJob(c *gin.Context) {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Cron.Remove(cron.EntryID(id))
	c.Status(http.StatusOK)
}

func NewJobController() *JobControllerImp {
	return &JobControllerImp{}
}
