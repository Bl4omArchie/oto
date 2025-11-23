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
	BinaryID		int				`gorm:"not null"`
	Binary			*Binary			`gorm:"foreignKey:BinaryID"`
	RequiresRoot	bool			`gorm:"not null"`
	RequiresValue	bool			`gorm:"not null"`
	ValueType		ValueType		`gorm:"not null"`
    ConflictsWith	[]Parameter		`gorm:"many2many:flag_conflicts;joinForeignKey:flag_id;joinReferences:conflict_id"`
    DependsOn		[]Parameter		`gorm:"many2many:flag_dependencies;joinForeignKey:flag_id;joinReferences:depends_on_id"`
}

// Newmodels.Parameter returns a new models.Parameter with a flag, description, the corresponding Binary ID, if the flag needs root access or a value and the value type 
// ValueType can be set to "" (None) 
func NewParameter(flag, description string, bin *Binary, requiresRoot, requiresValue bool, valueType ValueType, conflictsWith, dependsOn []Parameter) *Parameter {
	return &Parameter{
		Flag: flag,
		BinaryID: int(bin.ID),
		Binary: bin,
		Description: description,
		RequiresRoot: requiresRoot,
		RequiresValue: requiresValue,
		ValueType: valueType,
		ConflictsWith: conflictsWith,
		DependsOn: dependsOn,
	}
}

// FetchParameter returns the first parameter corresponding to the given column and given flag.
func FetchParameter(ctx context.Context, db *gorm.DB, column string, flag any) (*Parameter, error) {
	var param Parameter

	err := db.WithContext(ctx).
		Preload("Binary").
		Preload("ConflictsWith").
		Preload("DependsOn").
		Where(fmt.Sprintf("%s = ?", column), flag).
		First(&param).Error
	if err != nil {
		return nil, err
	}

	return &param, nil
}

// FetchParameter returns every parameters corresponding to the given column and given flag..
func FetchParameters(ctx context.Context, db *gorm.DB, column string, value any) ([]Parameter, error) {
	var params []Parameter

	err := db.WithContext(ctx).
		Preload("Binary").
		Preload("ConflictsWith").
		Preload("DependsOn").
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


func AllValueTypes() []ValueType {
	return []ValueType{Integer, String, Tuple, FilePath, Float, IPAddress, Port}
}
