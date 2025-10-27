package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/Bl4omArchie/oto/db"
	"github.com/Bl4omArchie/oto/models"
)

func CreateExecutable(name, version, binary, description string, c *gin.Context, cfg *models.Oto) {
	exec, err := models.NewExecutable(name, version, binary, description, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't create executable: ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}

func GetExecutables(c *gin.Context, cfg *models.Oto) {
	execs, err := db.GetTable[models.Executable](cfg.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get executables": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetExecutable(execId string, c *gin.Context, cfg *models.Oto) {
	exec, err := db.GetBy[models.Executable](cfg.Database, "exec_id", execId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get executable": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}

func GetExecutableInfo(execId string, c *gin.Context, cfg *models.Oto) {
	exec, err := db.GetBy[models.Executable](cfg.Database, "exec_id", execId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get executable info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}
