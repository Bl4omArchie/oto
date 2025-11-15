package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/Bl4omArchie/simple"
	"github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/oto/models"
	"github.com/go-playground/validator"

)

func CreateExecutable(c *gin.Context, cfg *oto.Oto) {
	var exec models.Executable

	if err := c.ShouldBindJSON(&exec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&exec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cfg.Database.Create(&exec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create executable", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exec)
}

func GetExecutables(c *gin.Context, cfg *oto.Oto) {
	execs, err := simple.GetRows[models.Executable](c, cfg.Database, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get executables": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetExecutable(execId string, c *gin.Context, cfg *oto.Oto) {
	exec, err := simple.GetRowBy[models.Executable](c, cfg.Database, "exec_id", execId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get executable": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}
