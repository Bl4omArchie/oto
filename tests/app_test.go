package test

import (
	"context"
	"testing"

	"github.com/Bl4omArchie/oto"
	"github.com/Bl4omArchie/oto/models"
	oto "github.com/Bl4omArchie/oto/pkg"
)


func TestOtoApp(t *testing.T) {
	oto, err := oto.NewInstanceOto()
	if err != nil {
		t.Fatalf("failed to get new instance of OTO : %v", err)
	}


}
