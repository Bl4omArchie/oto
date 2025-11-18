package models

import (
	"context"

	"github.com/Bl4omArchie/simple"
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
	BinID			int				`gorm:"not null"`
	Binary			*Binary			`gorm:"foreignKey:BinId"`
	RequiresRoot	bool			`gorm:"not null"`
	RequiresValue	bool			`gorm:"not null"`
	ValueType		ValueType		`gorm:"not null"`
    ConflictsWith	[]*Parameter	`gorm:"many2many:flag_conflicts;joinForeignKey:flag_id;joinReferences:conflict_id"`
    DependsOn		[]*Parameter	`gorm:"many2many:flag_dependencies;joinForeignKey:flag_id;joinReferences:depends_on_id"`
}

// Newmodels.Parameter returns a new models.Parameter with a flag, description, the corresponding Binary ID, if the flag needs root access or a value and the value type 
// ValueType can be set to "" (None) 
func NewParameter(flag, description string, bin *Binary, requiresRoot, requiresValue bool, valueType ValueType, conflictsWith, dependsOn []*Parameter) *Parameter {
	return &Parameter{
		Flag: flag,
		BinID: int(bin.ID),
		Binary: bin,
		Description: description,
		RequiresRoot: requiresRoot,
		RequiresValue: requiresValue,
		ValueType: valueType,
		ConflictsWith: conflictsWith,
		DependsOn: dependsOn,
	}
}

func FetchParameter(ctx context.Context, db *gorm.DB, field string, flag any) (*Parameter, error) {
	param, err := simple.GetRowBy[Parameter](ctx, db, field, flag)
	if err != nil {
		return nil, err
	}

	bin, err := FetchBinary(ctx, db, "ID", param.BinID)
	if err != nil {
		return nil, err
	}
	param.Binary = bin

	return param, nil
}

func FetchParameters(ctx context.Context, db *gorm.DB, flags []string) ([]*Parameter, error) {
	var result []*Parameter

	params, err := simple.GetRows[Parameter](ctx, db, "tag", tag)
		if err != nil {
			return nil, err
		}
		result = append(result, param)
	}
	return result, nil
}

func AllValueTypes() []ValueType {
	return []ValueType{Integer, String, Tuple, FilePath, Float, IPAddress, Port}
}
