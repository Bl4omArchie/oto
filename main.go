package main

import (
	"fmt"
	"context"

	"github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/oto/models"
)

func init_oto(cfg *oto.Instance) error {
	var ctx context.Context = context.Background()

	err := cfg.AddBinary("nmap", "7.98", "/usr/bin/nmap", "scanning tool")
	if err != nil {
		fmt.Println(err)
	}

	s, err := cfg.AddBinarySchema(ctx, "nmap - 7.98")
	if err != nil {
		return err
	}
	
	err = cfg.AddParameter(ctx, "nmap - 7.98", "-sL", "scan option for determine which host are online", false, false, models.String, []string{}, []string{}, s)
	if err != nil {
		return err
	}

	err = cfg.AddParameter(ctx, "nmap - 7.98", "-sT", "scan type", false, false, models.String, []string{"-sL"}, []string{}, s)
	if err != nil {
		return err
	}

	err = cfg.AddParameter(ctx, "nmap - 7.98", "-sK", "scan with -sT", false, false, models.String, []string{}, []string{"-sT"}, s)
	if err != nil {
		return err
	}

	err = cfg.AddParameter(ctx, "nmap - 7.98", "-T", "option depending on -sL", false, false, models.String, []string{"-sL", "-sT"}, []string{}, s)
	if err != nil {
		return err
	}

	err = cfg.AddCommand(ctx, "nmap - 7.98", "-sL", "determine which hosts are online", []string{"-sT", "-T"}, s)
	if err != nil {
		fmt.Println(err)
	}
	// cfg.AddJob(ctx, "nmap - 7.98", "-sL", "target gorm.io", map[string]string{"-sL": "185.199.111.153"})

	return nil
}


func main() {
	cfg, err := oto.NewInstanceOto("db/storage.db")
	if err != nil {
		fmt.Println(err)
	}

	init_oto(cfg)
}
