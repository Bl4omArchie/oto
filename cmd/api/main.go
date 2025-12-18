package main

import (
	"flag"
	"fmt"

	"github.com/Bl4omArchie/oto/api"
	oto "github.com/Bl4omArchie/oto/pkg"
)

func RunAPiServer(host, port string) error {
	cfg, err := oto.NewInstanceOto(".env")
	if err != nil {
		return err
	}

	r := api.SetupRouter(cfg)
	if err = r.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		return err
	}

	return nil
}

func main() {
	var host string
	var port string

	flag.StringVar(&host, "h", "0.0.0.0", "Specify host. Default is localhost")
	flag.StringVar(&port, "p", "8080", "Specify port number. Default is 8080")

	flag.Parse()

	RunAPiServer(host, port)
}
