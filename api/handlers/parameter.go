package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/Bl4omArchie/oto/db"
	"github.com/Bl4omArchie/oto/models"
)


func CreateParameter(execId, flag, description string, reqRoot, reqValue bool, valueType models.ValueType, dependsOn []string, conflictWith []string, c *gin.Context, oto *models.Oto) {
	exec, err := models.NewParameter(execId, flag, description, reqRoot, reqValue, valueType, dependsOn, conflictWith, oto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't create parameter": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}

func GetParameters(c *gin.Context, oto *models.Oto) {
	execs, err := db.GetTable[models.Parameter](oto.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get parameters": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetParameter(execId string, paramName string, c *gin.Context, oto *models.Oto) {
	param, err := db.GetBy[models.Parameter](oto.Database, execId, paramName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get parameter": err.Error()})
		return
	}
	c.JSON(http.StatusOK, param)
}

func GetValueTypes(c *gin.Context) {
    c.JSON(http.StatusOK, []string{"string", "int", "bool"})
}
