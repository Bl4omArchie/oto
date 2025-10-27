package api

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/oto/models"
	"github.com/Bl4omArchie/oto/api/handlers"
)

func SetupRouter(cfg *pkg.Oto) *gin.Engine {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r.GET("/execs", func(c *gin.Context) {
		handlers.GetExecutables(c, cfg)
	})

    r.GET("/execs/:execId", func(c *gin.Context) {
		value := c.Param("execId")
		handlers.GetExecutable(value, c, cfg)
	})

    r.GET("/params/:execId", func(c *gin.Context) {
		handlers.GetParameters(c, cfg)
	})

    r.GET("/cmds/:execId", func(c *gin.Context) {
		handlers.GetCommands(c, cfg)
	})

    r.GET("/params/:execId/:name", func(c *gin.Context) {
		key := c.Param("execId")
		value := c.Param("name")
		handlers.GetParameter(key, value, c, cfg)
	})

	r.GET("/cmds/:execId/:name", func(c *gin.Context) {
		key := c.Param("execId")
		value := c.Param("name")
		handlers.GetCommand(key, value, c, cfg)
	})

	r.GET("/valuetypes", func(c *gin.Context) {
		handlers.GetValueTypes(c)
	})

    r.POST("/execs", func(c *gin.Context) {
		var exec models.Executable
		if err := c.ShouldBindJSON(&exec); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cfg.Database.Create(exec)
	})

    r.POST("/params", func(c *gin.Context) {
		var param models.Parameter
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cfg.Database.Create(param)
	})

    r.POST("/cmds", func(c *gin.Context) {
		var cmd models.Command
		if err := c.ShouldBindJSON(&cmd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cfg.Database.Create(cmd)
	})

    return r
}
