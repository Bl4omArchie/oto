package handlers

import (
	"net/http"

	"github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/simple"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func CreateExecutable(c *gin.Context, cfg *oto.Instance) {
	var executable models.Executable

	if err := c.ShouldBindJSON(&executable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&executable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cfg.Database.Create(&executable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Executable", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, executable)
}

func GetExecutables(c *gin.Context, cfg *oto.Instance) error {
	executables, err := simple.GetTable[models.Executable](c, cfg.Database, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get Executables": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, executables)
	return nil
}

func GetExecutable(binName string, c *gin.Context, cfg *oto.Instance) error {
	executable, err := simple.GetRowBy[models.Executable](c, cfg.Database, "tag", binName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get Executable": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, executable)
	return nil
}
