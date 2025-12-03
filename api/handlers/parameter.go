package handlers

import (
	"net/http"

	"github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/simple"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func CreateParameter(c *gin.Context, cfg *oto.Instance) {
	var param models.Parameter

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cfg.Database.Create(&param).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create parameter", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, param)
}

func GetParameters(c *gin.Context, oto *oto.Instance) {
	execs, err := simple.GetTable[models.Parameter](c, oto.Database, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get parameters": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetParameter(binTag string, paramName string, c *gin.Context, oto *oto.Instance) {
	param, err := simple.GetRowBy[models.Parameter](c, oto.Database, binTag, paramName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get parameter": err.Error()})
		return
	}
	c.JSON(http.StatusOK, param)
}

func GetValueTypes(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"string", "int", "bool"})
}
