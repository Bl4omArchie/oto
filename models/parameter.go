package models


import (
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

type FlagID string
type Parameter struct {
	gorm.Model
	Flag          string         `gorm:"unique;not null" json:"flag" validate:"required"`
	Description   string         `gorm:"type:text" json:"description" validate:"omitempty"`
	ExecID        int         	 `gorm:"not null" json:"exec_id" validate:"required"`
	Exec		  Executable
	RequiresRoot  bool           `gorm:"not null"  json:"requires_root" validate:"required"`
	RequiresValue bool           `gorm:"not null"  json:"requires_value" validate:"required"`
	ValueType     ValueType      `gorm:"not null"  json:"value_type" validate:"required"`
    ConflictsWith []*Parameter 	 `gorm:"many2many:flag_conflicts;joinForeignKey:flag_id;joinReferences:conflict_id" json:"conflict_with" validate:"omitempty"`
    DependsOn     []*Parameter 	 `gorm:"many2many:flag_dependencies;joinForeignKey:flag_id;joinReferences:depends_on_id"  json:"depends_on" validate:"omitempty"`
}

// Newmodels.Parameter returns a new models.Parameter with a flag, description, the corresponding executable ID, if the flag needs root access or a value and the value type 
// ValueType can be set to "" (None) 
func NewParameter(flag, description string, exec *Executable, requiresRoot, requiresValue bool, valueType ValueType, conflictsWith, dependsOn []*Parameter) *Parameter {
	return &Parameter{
		Flag: flag,
		ExecID: int(exec.ID),
		Exec: *exec,
		Description: description,
		RequiresRoot: requiresRoot,
		RequiresValue: requiresValue,
		ValueType: valueType,
		ConflictsWith: conflictsWith,
		DependsOn: dependsOn,
	}
}


func AllValueTypes() []ValueType {
	return []ValueType{Integer, String, Tuple, FilePath, Float, IPAddress, Port}
}
