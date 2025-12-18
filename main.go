package main

import (
	"context"
	"fmt"

	_ "github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"

	_ "ariga.io/atlas-provider-gorm/gormschema"
)

func full_data() error {
	var ctx context.Context = context.Background()

	instance, err := oto.NewInstanceOto(".env")
	if err != nil {
		return err
	}

	err = instance.AddExecutable("nmap", "7.98", "/usr/bin/nmap", "scanning tool")
	if err != nil {
		fmt.Println(err)
	}

	err = instance.AddExecutable("openssl", "3.5.3", "/usr/bin/openssl", "cryptographic tool")
	if err != nil {
		fmt.Println(err)
	}

	err = instance.AddExecutable("masscan", "1.3.9", "/usr/bin/masscan", "scanning tool")
	if err != nil {
		fmt.Println(err)
	}

	s1, err := instance.AddExecutableSchema(ctx, "nmap - 7.98")
	if err != nil {
		return err
	}
	s2, err := instance.AddExecutableSchema(ctx, "openssl - 3.5.3")
	if err != nil {
		return err
	}
	s3, err := instance.AddExecutableSchema(ctx, "masscan - 1.3.9")
	if err != nil {
		return err
	}

	err = instance.ImportParameters(ctx, "data/nmap.json", s1)
	if err != nil {
		return err
	}

	err = instance.ImportParameters(ctx, "data/openssl.json", s2)
	if err != nil {
		return err
	}

	err = instance.ImportParameters(ctx, "data/masscan.json", s3)
	if err != nil {
		return err
	}

	err = instance.AddCommand(ctx, "openssl - 3.5.3", "GenRSA", "Generate an rsa keypair", []string{"genpkey", "-algorithm", "-pkeyopt", "-out"}, s2)
	if err != nil {
		return err
	}

	if err := instance.AddJob(ctx, "GenRSA", "GenRSA-2048", map[string]string{"genpkey": "", "-algorithm": "RSA", "-pkeyopt": "rsa_keygen_bits:2048", "-out": "key.pem"}); err != nil {
		return err
	}

	out, err := instance.RunJobDemo(ctx, "GenRSA-2048")
	if err != nil {
		return err
	}

	fmt.Println(out.Stderr, out.Stdout)

	return nil
}

func launch_demo() error {
	var ctx context.Context = context.Background()

	instance, err := oto.NewInstanceOto(".env")
	if err != nil {
		return err
	}

	err = instance.AddExecutable("openssl", "3.5.3", "/usr/bin/openssl", "cryptographic tool")
	if err != nil {
		fmt.Println(err)
	}

	s2, err := instance.AddExecutableSchema(ctx, "openssl - 3.5.3")
	if err != nil {
		return err
	}
	err = instance.ImportParameters(ctx, "data/openssl.json", s2)
	if err != nil {
		return err
	}

	err = instance.AddCommand(ctx, "openssl - 3.5.3", "GenRSA", "Generate an rsa keypair", []string{"genpkey", "-algorithm", "-pkeyopt", "-out"}, s2)
	if err != nil {
		return err
	}

	if err := instance.AddJob(ctx, "GenRSA", "GenRSA-2048", map[string]string{"genpkey": "", "-algorithm": "RSA", "-pkeyopt": "rsa_keygen_bits:2048", "-out": "key.pem"}); err != nil {
		return err
	}

	out, err := instance.RunJobDemo(ctx, "GenRSA-2048")
	if err != nil {
		return err
	}

	fmt.Println(out.Stderr, out.Stdout)
	return nil
}

func test() error {
	var ctx context.Context = context.Background()

	instance, err := oto.NewInstanceOto(".env")
	if err != nil {
		return err
	}

	instance.AddRoutine(ctx, "test_routine", "")

	return nil
}

func main() {
	fmt.Println(test())
}
