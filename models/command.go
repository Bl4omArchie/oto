package models

import (
	"fmt"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Command struct {
	gorm.Model
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string         `gorm:"unique;not null" json:"name" validate:"required"`
	ExecID         string         `gorm:"not null" json:"exec_id" validate:"required"`
	RequiresRoot   bool           `gorm:"not null" json:"requires_root" validate:"required"`
	ParameterFlags datatypes.JSON `gorm:"type:json" json:"parameter_flags" validate:"required"`
	Description    string         `gorm:"type:text" json:"description" validate:"omitempty"`
}

func NewCommand(execId, cmdName, description string, flags []string) (*Command, error) {
	flagsBytes, err := SliceToJson(flags)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert flags to json")
	}

	return &Command{
		Name: cmdName,
		ExecID: execId,
		Description: description,
		ParameterFlags: flagsBytes,
	}, nil
}

func SliceToJson(v []string) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}
