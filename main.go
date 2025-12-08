package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"

	_ "ariga.io/atlas-provider-gorm/gormschema"
)

func fill_database() error {
	var ctx context.Context = context.Background()

	instance, err := oto.NewInstanceOto(".env")
	if err != nil {
		return err
	}

	// err = instance.AddExecutable("nmap", "7.98", "/usr/exec/nmap", "scanning tool")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = instance.AddExecutable("openssl", "3.5.3", "/usr/exec/openssl", "cryptographic tool")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = instance.AddExecutable("masscan", "1.3.9", "/usr/exec/masscan", "scanning tool")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	s, err := instance.AddExecutableSchema(ctx, "nmap - 7.98")
	if err != nil {
		return err
	}
	s, err = instance.AddExecutableSchema(ctx, "openssl - 3.5.3")
	if err != nil {
		return err
	}
	s, err = instance.AddExecutableSchema(ctx, "masscan - 1.3.9")
	if err != nil {
		return err
	}

	// err = instance.ImportParameters(ctx, "data/nmap.json", s)
	// if err != nil {
	// 	return nil
	// }

	err = instance.ImportParameters(ctx, "data/openssl.json", s)
	if err != nil {
		return nil
	}

	err = instance.ImportParameters(ctx, "data/masscan.json", s)
	if err != nil {
		return nil
	}

	// err = instance.AddCommand(ctx, "nmap - 7.98", "reco", "determine which hosts are online", []string{"-PS"}, s)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// arg1, err := models.FetchParameter(ctx, instance.Database, "flag", "-PS")
	// instance.AddJob(ctx, "reco", "target gorm.io", map[*models.Parameter]string{arg1: "185.199.111.153"})

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
	fmt.Println(fill_database())
}
