package handlers

import (
	"net/http"

	"github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/simple"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func CreateBinary(c *gin.Context, cfg *oto.Instance) {
	var binary models.Binary

	if err := c.ShouldBindJSON(&binary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&binary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cfg.Database.Create(&binary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Binary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, binary)
}

func GetBinaries(c *gin.Context, cfg *oto.Instance) {
	binarys, err := simple.GetTable[models.Binary](c, cfg.Database, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get Binaries": err.Error()})
		return
	}
	c.JSON(http.StatusOK, binarys)
}

func GetBinary(binTag string, c *gin.Context, cfg *oto.Instance) {
	binary, err := simple.GetRowBy[models.Binary](c, cfg.Database, "tag", binTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get Binary": err.Error()})
		return
	}
	c.JSON(http.StatusOK, binary)
}
