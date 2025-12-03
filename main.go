package main

import (
	"fmt"
	"time"
	"context"

	"github.com/Bl4omArchie/oto/pkg"
	"github.com/Bl4omArchie/oto/models"
)

func fill_database() error {
	var ctx context.Context = context.Background()

	instance, err := oto.NewInstanceOto(".env")
	if err != nil {
		return err
	}

	err = instance.AddBinary("nmap", "7.98", "/usr/bin/nmap", "scanning tool")
	if err != nil {
		fmt.Println(err)
	}

	s, err := instance.AddBinarySchema(ctx, "nmap - 7.98")
	if err != nil {
		return err
	}

	err = instance.ImportParameters(ctx, "models/nmap.json", s)
	if err != nil {
		return nil
	}

	err = instance.AddCommand(ctx, "nmap - 7.98", "reco", "determine which hosts are online", []string{"-PS"}, s)
	if err != nil {
		fmt.Println(err)
	}
	
	arg1, err := models.FetchParameter(ctx, instance.Database, "flag", "-PS")
	instance.AddJob(ctx, "reco", "target gorm.io", map[*models.Parameter]string{arg1: "185.199.111.153"})

	return nil
}

func temporal_test() error {
	var ctx context.Context = context.Background()

	instance, err := oto.NewInstanceOto(".env")
	if err != nil {
		return err
	}

	go func() error {
		if err := instance.StartWorker("oto-tasks"); err != nil {
			return err
		}
		return nil
	}()

	time.Sleep(500 * time.Millisecond)

	result, err := instance.RunJobWorkflow(ctx, "target gorm.io")
	if err != nil {
		return err
	}
	fmt.Println("Workflow finished:", result.Stdout)
	
	return nil
}

func main() {
	fmt.Println(temporal_test())
}
