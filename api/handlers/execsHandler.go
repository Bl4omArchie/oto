package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/Bl4omArchie/oto/db"
	"github.com/Bl4omArchie/oto/pkg"
)


func GetExecutables(c *gin.Context, oto *pkg.Oto) {
	execs, err := db.GetTable[pkg.Executable](oto.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}


func GetCommands(c *gin.Context, oto *pkg.Oto) {
	execs, err := db.GetTable[pkg.Command](oto.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}


func GetParameters(c *gin.Context, oto *pkg.Oto) {
	execs, err := db.GetTable[pkg.Parameter](oto.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execs)
}

func GetExecutable(key string, value string, c *gin.Context, oto *pkg.Oto) {
	exec, err := db.GetBy[pkg.Executable](oto.Database, key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exec)
}

func GetParameter(key string, value string, c *gin.Context, oto *pkg.Oto) {
	param, err := db.GetBy[pkg.Parameter](oto.Database, key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, param)
}

func GetCommand(key string, value string, c *gin.Context, oto *pkg.Oto) {
	cmd, err := db.GetBy[pkg.Command](oto.Database, key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cmd)
}
