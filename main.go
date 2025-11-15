package main

import (
	"fmt"
	"context"

	"github.com/Bl4omArchie/oto/pkg"
	_ "github.com/Bl4omArchie/oto/models"
)

func main() {
	oto, err := oto.NewInstanceOto("db/storage.db")
	if err != nil {
		fmt.Println(err)
	}

	var ctx context.Context = context.Background()
	// oto.AddExecutable("nmap", "7.98", "/usr/bin/nmap", "scanning tool")
	// err =  oto.AddParameter(ctx, "nmap - 7.98", "-sL", "scan option for determine which hostt are online", false, false, models.String, []string{}, []string{})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = oto.AddCommand(ctx, "nmap - 7.98", "-sL", "determine which hosts are online", []string{"-sL"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// oto.AddJobCommand(ctx, "nmap - 7.98", "-sL", "target gorm.io", map[string]string{"-sL": "185.199.111.153"})

	output, err := oto.RunJobCommand(ctx, "target gorm.io")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)

}
