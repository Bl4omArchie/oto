package models

import (
	"fmt"
	"context"

	"gorm.io/gorm"
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
	Flag			string			`gorm:"unique;not null"`
	Description		string			`gorm:"type:text"`
	BinaryTag		string			`gorm:"not null"`
	Binary			*Binary			`gorm:"foreignKey:BinaryTag"`
	RequiresRoot	bool			`gorm:"not null"`
	RequiresValue	bool			`gorm:"not null"`
	ValueType		ValueType		`gorm:"not null"`
	Require			[]Parameter		`gorm:"many2many:flag_dependencies;joinForeignKey:flag_id;joinReferences:requires_id"`
    Interfer		[]Parameter		`gorm:"many2many:flag_conflicts;joinForeignKey:flag_id;joinReferences:interfer_id"`
}

type ParameterRaw struct {
    Flag          string   `json:"flag"`
    Description   string   `json:"description"`
    BinaryTag     string   `json:"binary_tag"`
    RequiresRoot  bool     `json:"requires_root"`
    RequiresValue bool     `json:"requires_value"`
    ValueType     ValueType`json:"value_type"`
    RequireIDs    []string `json:"require_ids"`
    InterferIDs   []string `json:"interfer_ids"`
}

// Newmodels.Parameter returns a new models.Parameter with a flag, description, the corresponding Binary ID, if the flag needs root access or a value and the value type 
// ValueType can be set to "" (None) 
func NewParameter(flag, description string, bin *Binary, requiresRoot, requiresValue bool, valueType ValueType, require, interfer []Parameter) *Parameter {
	return &Parameter{
		Flag: flag,
		BinaryTag: bin.Tag,
		Binary: bin,
		Description: description,
		RequiresRoot: requiresRoot,
		RequiresValue: requiresValue,
		ValueType: valueType,
		Interfer: interfer,
		Require: require,
	}
}

func NewParameterRaw(flag, description, binaryTag string, requiresRoot, requiresValue bool, valueType ValueType, require, interfer []string) *ParameterRaw {
	return &ParameterRaw{
		Flag: flag,
		Description: description,
		BinaryTag: binaryTag,
		RequiresRoot: requiresRoot,
		RequiresValue: requiresValue,
		ValueType: valueType,
		RequireIDs: require,
		InterferIDs: interfer,
	}
}

// FetchParameter returns the first parameter corresponding to the given column and value.
func FetchParameter(ctx context.Context, db *gorm.DB, column string, value any) (*Parameter, error) {
	var param Parameter

	err := db.WithContext(ctx).
		Preload("Binary").
		Preload("Interfer").
		Preload("Require").
		Where(fmt.Sprintf("%s = ?", column), value).
		First(&param).Error
	if err != nil {
		return nil, err
	}

	return &param, nil
}

// FetchParameter returns every parameters corresponding to the given column and value.
func FetchParameters(ctx context.Context, db *gorm.DB, column string, value any) ([]Parameter, error) {
	var params []Parameter

	err := db.WithContext(ctx).
		Preload("Binary").
		Preload("Interfer").
		Preload("Require").
		Where(fmt.Sprintf("%s = ?", column), value).
		Find(&params).Error
	if err != nil {
		return nil, err
	}

	return params, nil
}

// FetchFlagParameters returns parameters corresponding to the given flags
func FetchFlagParameters(ctx context.Context, db *gorm.DB, column string, flags []string) ([]Parameter, error) {
	var result []Parameter

	for _, flag := range flags {
		param, err := FetchParameter(ctx, db, column, flag)
		if err != nil {
			return nil, err
		}
		result = append(result, *param)
	}
	return result, nil
}

// AllValueTypes list every supported type for a parameter value
func AllValueTypes() []ValueType {
	return []ValueType{Integer, String, Tuple, FilePath, Float, IPAddress, Port}
}
