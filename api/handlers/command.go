package handlers

import (
	"net/http"

	"github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/simple"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func CreateCommand(c *gin.Context, cfg *oto.Instance) {
	var cmd models.Command

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cfg.Database.Create(&cmd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create command", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cmd)
}

func GetCommands(c *gin.Context, oto *oto.Instance) {
	cmds, err := simple.GetTable[models.Command](c, oto.Database, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get commands": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cmds)
}

func GetCommand(execTag string, cmdName string, c *gin.Context, oto *oto.Instance) {
	cmd, err := simple.GetRowBy[models.Command](c, oto.Database, "name", cmdName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get command": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cmd)
}
