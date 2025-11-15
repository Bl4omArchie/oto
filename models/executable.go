package models

import (
	"fmt"

	"gorm.io/gorm"
)


type Executable struct {
	gorm.Model
	ExecId		string `gorm:"unique; not null" json:"exec_id" validate:"required"`
	Name        string `gorm:"not null type:string" json:"name" validate:"required"`
	Version     string `gorm:"not null type:string" json:"version" validate:"required"`
	Binary      string `gorm:"not null type:string" json:"binary" validate:"required"`
	Description string `gorm:"not null type:string" json:"description" validate:"required"`
}

func NewExecutable(name, version, binaryPath, description string) *Executable {
	execId := GetExecId(name, version)
	return &Executable {
		ExecId: execId,
		Name: name,
		Version: version,
		Binary: binaryPath,
		Description: description,
	}
}

func GetExecId(name string, version string) string {
	return fmt.Sprintf("%s - %s", name, version)
}
