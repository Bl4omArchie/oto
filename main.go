package main

import (
    "fmt"
    "github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/oto/api"
)

func StartApi(cfg *pkg.Oto) {
    r := api.SetupRouter(cfg)
    r.Run(":8080")
}

func main() {
	cfg, err := pkg.OpenOto("db/oto_storage.db")
	if err != nil {
		fmt.Println(err)
	}

	StartApi(cfg)

    //oto.NewExecutable("nmap", "7.94SVN", "/usr/bin/nmap", "port scanning tool", config)
    //oto.NewParameter("nmap - 7.94SVN", "-sn", "basic tcp scan", false, false, oto.String, []string{}, []string{}, config)
    //oto.NewCommand("nmap - 7.94SVN", "tcp scan", "basic tcp scan", []string{"-sn"}, config)

	//oto.UpdateTable[oto.Parameter]("flag", "-sL", "exec_id", "nmap - 7.94SVN", config)

}
