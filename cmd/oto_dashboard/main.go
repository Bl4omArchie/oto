package main

import (
	"fmt"
	"flag"

    "github.com/gin-gonic/gin"
)

func RunWebDashboard(host, port string) {
    r := gin.Default()

    r.Static("/web", "./web")

    r.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.html", nil)
    })

    r.GET("/executables", func(c *gin.Context) {
        c.HTML(200, "executables.html", nil)
    })

    r.GET("/executables/info", func(c *gin.Context) {
        c.HTML(200, "exec-info.html", nil)
    })


    r.Run(fmt.Sprintf("%s:%s", host, port))
}


func main() {
    var host string    
    var port string      
 
    flag.StringVar(&host, "h", "localhost", "Specify host. Default is localhost")
    flag.StringVar(&port, "p", "9090", "Specify port number. Default is 9090")

	flag.Parse()
 
	RunWebDashboard(host, port)
}
