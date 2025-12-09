package main


import (
	"fmt"
	"flag"

	"github.com/Bl4omArchie/oto/api"
	"github.com/Bl4omArchie/oto/pkg"
)

func RunAPiServer(host, port string) {
	cfg, err := oto.NewInstanceOto(".env")
	if err != nil {
		fmt.Println(err)
	}

    r := api.SetupRouter(cfg)
    r.Run(fmt.Sprintf("%s:%s", host, port))
}


func main() {
    var host string    
    var port string      
 
    flag.StringVar(&host, "h", "0.0.0.0", "Specify host. Default is localhost")
    flag.StringVar(&port, "p", "8080", "Specify port number. Default is 8080")
	
	flag.Parse()

	RunAPiServer(host, port)
}
