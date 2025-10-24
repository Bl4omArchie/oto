package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/Bl4omArchie/oto/db"
	"github.com/Bl4omArchie/oto/pkg"
)


func CreateParameter(execId, flag, description string, reqRoot, reqValue bool, valueType pkg.ValueType, dependsOn []string, conflictWith []string, c *gin.Context, oto *pkg.Oto) {
	exec, err := pkg.NewParameter(execId, flag, description, reqRoot, reqValue, valueType, dependsOn, conflictWith, oto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't create parameter": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}

func GetParameters(c *gin.Context, oto *pkg.Oto) {
	execs, err := db.GetTable[pkg.Parameter](oto.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get parameters": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetParameter(execId string, paramName string, c *gin.Context, oto *pkg.Oto) {
	param, err := db.GetBy[pkg.Parameter](oto.Database, execId, paramName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error, couldn't get parameter": err.Error()})
		return
	}
	c.JSON(http.StatusOK, param)
}

func GetValueTypes(c *gin.Context) {
    c.JSON(http.StatusOK, []string{"string", "int", "bool"})
}
