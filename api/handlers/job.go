package handlers

import (
	"net/http"

	"github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/simple"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func CreateJob(c *gin.Context, cfg *oto.Instance) {
	var job models.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cfg.Database.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context, oto *oto.Instance) {
	execs, err := simple.GetTable[models.Job](c, oto.Database, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get commands": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetJob(jobName string, c *gin.Context, oto *oto.Instance) {
	jobCmd, err := simple.GetRowBy[models.Job](c, oto.Database, "name", jobName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get command": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobCmd)
}
