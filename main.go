package main

import (
	"fmt"
	"context"

	"github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/oto/models"
)

func init_oto(cfg *oto.Config) {
	var ctx context.Context = context.Background()

	cfg.AddExecutable("nmap", "7.98", "/usr/bin/nmap", "scanning tool")
	err :=  cfg.AddParameter(ctx, "nmap - 7.98", "-sL", "scan option for determine which hostt are online", false, false, models.String, []string{}, []string{})
	if err != nil {
		fmt.Println(err)
	}

	err = cfg.AddCommand(ctx, "nmap - 7.98", "-sL", "determine which hosts are online", []string{"-sL"})
	if err != nil {
		fmt.Println(err)
	}
	cfg.AddJobCommand(ctx, "nmap - 7.98", "-sL", "target gorm.io", map[string]string{"-sL": "185.199.111.153"})

	output, err := cfg.RunJobCommand(ctx, "target gorm.io")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func flag_matching(cfg *oto.Config) {
	schema := oto.NewSchema([]models.FlagID{"a", "b", "c", "d", "e"})
	schema.Require("a", "b")
	schema.Require("b", "c")

	fmt.Println(schema.ValidateSchema())

}

func main() {
	cfg, err := oto.NewInstanceOto("db/storage.db")
	if err != nil {
		fmt.Println(err)
	}

	flag_matching(cfg)
}
