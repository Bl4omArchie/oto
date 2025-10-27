package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/Bl4omArchie/oto/db"
	"github.com/Bl4omArchie/oto/models"
)


func CreateCommand(execId, cmdName, description string, flags []string, c *gin.Context, oto *models.Oto) {
	exec, err := models.NewCommand(execId, cmdName, description, flags, oto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't create command": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}

func GetCommands(c *gin.Context, oto *models.Oto) {
	execs, err := db.GetTable[models.Command](oto.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get commands": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetCommand(execId string, cmdName string, c *gin.Context, oto *models.Oto) {
	cmd, err := db.GetBy[models.Command](oto.Database, execId, cmdName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get command": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cmd)
}
