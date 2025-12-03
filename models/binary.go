package models

import (
	"context"
	"fmt"

	"github.com/Bl4omArchie/simple"
	"gorm.io/gorm"
)

type Binary struct {
	gorm.Model
	Tag         string `gorm:"unique; not null type:string"`
	Name        string `gorm:"not null type:string"`
	Version     string `gorm:"not null type:string"`
	Path        string `gorm:"not null type:string"`
	Description string `gorm:"not null type:string"`
}

func NewBinary(name, version, path, description string) *Binary {
	tag := GetTag(name, version)
	return &Binary{
		Tag:         tag,
		Name:        name,
		Version:     version,
		Path:        path,
		Description: description,
	}
}

func FetchBinary(ctx context.Context, db *gorm.DB, column string, tag any) (*Binary, error) {
	return simple.GetRowBy[Binary](ctx, db, column, tag)
}

func GetTag(name string, version string) string {
	return fmt.Sprintf("%s - %s", name, version)
}
