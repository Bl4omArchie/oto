package models



import (
	"context"
	"fmt"

	"github.com/Bl4omArchie/simple"
	"gorm.io/gorm"
)


// A executable is a binary with the execution rigth
//
// In order to create a new executable, you must specify its path and version
//
// You can find most of your installed executable like openssl, nmap etc in /usr/bin
//
// The executable Tag is the name and version combined into string (i.e : oppenssl - v3.5.3) in order to avoid confusion between different versions for a same executable
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
