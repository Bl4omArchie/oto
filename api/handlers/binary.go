package handlers

import (
	"net/http"

	"github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/simple"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func CreateBinary(c *gin.Context, cfg *oto.Config) {
	var exec models.Binary

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Binary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exec)
}

func GetBinarys(c *gin.Context, cfg *oto.Config) {
	execs, err := simple.GetRows[models.Binary](c, cfg.Database, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get Binarys": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetBinary(execId string, c *gin.Context, cfg *oto.Config) {
	exec, err := simple.GetRowBy[models.Binary](c, cfg.Database, "tag", execId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get Binary": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}
