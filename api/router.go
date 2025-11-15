package api

import (
	"time"

	"github.com/Bl4omArchie/oto/api/handlers"
	"github.com/Bl4omArchie/oto/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *oto.Oto) *gin.Engine {
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

    r.GET("/jobs", func(c *gin.Context) {
		handlers.GetJobCommands(c, cfg)
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

	r.GET("/jobs/:name", func(c *gin.Context) {
		value := c.Param("name")
		handlers.GetJobCommand(value, c, cfg)
	})

	r.GET("/valuetypes", func(c *gin.Context) {
		handlers.GetValueTypes(c)
	})

    r.POST("/execs", func(c *gin.Context) {
		handlers.CreateExecutable(c, cfg)
	})

    r.POST("/params", func(c *gin.Context) {
		handlers.CreateCommand(c, cfg)
	})

    r.POST("/cmds", func(c *gin.Context) {
		handlers.CreateParameter(c, cfg)
	})

    r.POST("/jobs", func(c *gin.Context) {
		handlers.CreateJobCommand(c, cfg)
	})
    return r
}
