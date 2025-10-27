package models


import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/datatypes"
)


type ValueType string
const (
	Integer   ValueType = "int"
	String    ValueType = "string"
	Tuple     ValueType = "tuple"
	FilePath  ValueType = "path"
	Float     ValueType = "float"
	IPAddress ValueType = "ip"
	Port      ValueType = "port"
	None	  ValueType = ""
)

type Parameter struct {
	gorm.Model
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Flag          string         `gorm:"unique;not null" json:"flag" validate:"required"`
	ExecID        string         `gorm:"not null"  json:"exec_id" validate:"required"`
	RequiresRoot  bool           `gorm:"not null"  json:"requires_root" validate:"required"`
	RequiresValue bool           `gorm:"not null"  json:"requires_value" validate:"required"`
	ValueType     ValueType      `gorm:"not null"  json:"value_type"`
	DependsOn     datatypes.JSON `gorm:"type:json" json:"depends_on" validate:"omitempty"`
	ConflictWith  datatypes.JSON `gorm:"type:json" json:"conflict_with" validate:"omitempty"`
	Description   string         `gorm:"type:text" json:"description" validate:"omitempty"`
}

func NewParameter(execId, flag, description string, reqRoot, reqValue bool, valueType ValueType, dependsOn []string, conflictWith []string) (*Parameter, error) {
	dependsOnBytes, err := SliceToJson(dependsOn)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert dependsOn to json")
	}
	conflictWithBytes, err := SliceToJson(conflictWith)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert conflictWith to json")
	}
	
	return &Parameter{
		Flag: flag,
		ExecID: execId,
		Description: description,
		RequiresRoot: reqRoot,
		RequiresValue: reqValue,
		ValueType: valueType,
		DependsOn: dependsOnBytes,
		ConflictWith: conflictWithBytes,
	}, nil
}

func AllValueTypes() []ValueType {
	return []ValueType{Integer, String, Tuple, FilePath, Float, IPAddress, Port}
}
