package oto

import (
	"context"
	"testing"

	"github.com/Bl4omArchie/oto/models"
)

func TestFillDatabase(t *testing.T) {
	instance, err := NewInstanceOto("../test.env")
	if err != nil {
		t.Fatalf("failed to get new instance of OTO : %v", err)
	}

	var ctx context.Context = context.Background()

	err = instance.AddExecutable("nmap", "7.98", "/usr/exec/nmap", "scanning tool")
	if err != nil {
		t.Fatalf("%v", err)
	}

	s, err := instance.AddExecutableSchema(ctx, "nmap - 7.98")
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddParameter(ctx, "nmap - 7.98", "-sL", "scan option for determine which host are online", false, false, models.String, []string{}, []string{}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddParameter(ctx, "nmap - 7.98", "-sT", "scan type", false, false, models.String, []string{"-sL"}, []string{}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddParameter(ctx, "nmap - 7.98", "-sK", "scan with -sT", false, false, models.String, []string{}, []string{"-sT"}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddParameter(ctx, "nmap - 7.98", "-T", "option depending on -sL", false, false, models.String, []string{"-sL", "-sT"}, []string{}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddCommand(ctx, "nmap - 7.98", "-sL", "determine which hosts are online", []string{"-sT", "-T"}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIncorrectSchema(t *testing.T) {
	instance, err := NewInstanceOto("../test.env")
	if err != nil {
		t.Fatalf("failed to get new instance of OTO : %v", err)
	}

	var ctx context.Context = context.Background()

	s, err := instance.AddExecutableSchema(ctx, "nmap - 7.98")
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddParameter(ctx, "nmap - 7.98", "c", "scan with -sT", false, false, models.String, []string{}, []string{"a"}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddParameter(ctx, "nmap - 7.98", "b", "scan type", false, false, models.String, []string{"c"}, []string{}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = instance.AddParameter(ctx, "nmap - 7.98", "a", "scan option for determine which host are online", false, false, models.String, []string{"b"}, []string{}, s)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := s.ValidateSchema(); err == nil {
		t.Fatalf("Schema validation didn't fail even though it should have")
	}
}
