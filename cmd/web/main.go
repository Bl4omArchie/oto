package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func RunWebDashboard(host, port, path string) {
    fs := http.FileServer(http.Dir(path))
    http.Handle("/", fs)

    log.Printf("Server starting at http://localhost:%s\n", port)
    err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
    if err != nil {
        log.Fatal(err)
    }
}


func main() {
    var host string    
    var port string
	var path string  
 
    flag.StringVar(&host, "h", "localhost", "Specify host. Default is localhost")
    flag.StringVar(&port, "p", "9090", "Specify port number. Default is 9090")
	flag.StringVar(&path, "path", "../../web/", "Specify the path of the OTO dashboard")

	flag.Parse()
 
	RunWebDashboard(host, port, path)
}
