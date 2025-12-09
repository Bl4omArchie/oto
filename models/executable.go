package models

import (
	"context"
	"fmt"

	"github.com/Bl4omArchie/simple"
	"gorm.io/gorm"
)

type Executable struct {
	gorm.Model
	Tag         string `gorm:"unique; not null type:string"`
	Name        string `gorm:"not null type:string"`
	Version     string `gorm:"not null type:string"`
	Path        string `gorm:"not null type:string"`
	Description string `gorm:"not null type:string"`
}

func NewExecutable(name, version, path, description string) *Executable {
	tag := GetTag(name, version)
	return &Executable{
		Tag:         tag,
		Name:        name,
		Version:     version,
		Path:        path,
		Description: description,
	}
}

func FetchExecutable(ctx context.Context, db *gorm.DB, column string, tag any) (*Executable, error) {
	return simple.GetRowBy[Executable](ctx, db, column, tag)
}

func GetTag(name string, version string) string {
	return fmt.Sprintf("%s - %s", name, version)
}
