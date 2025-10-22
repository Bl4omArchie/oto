package api

import (
	"github.com/gin-gonic/gin"
	"github.com/Bl4omArchie/oto/api/handlers"
	"github.com/Bl4omArchie/oto/pkg"
)

func SetupRouter(cfg *pkg.Oto) *gin.Engine {
    r := gin.Default()

    r.GET("/execs", handlers.GetExecutables)
	r.GET("/execs/:execId", handlers.GetExecutable)
	r.GET("/execs/:execId/params", handlers.GetParameters)
	r.GET("/execs/:execId/cmds", handlers.GetCommands)
	r.GET("/ping", handlers.Ping)
	r.GET("/ping", handlers.Ping)
	r.GET("/ping", handlers.Ping)

    return r
}
