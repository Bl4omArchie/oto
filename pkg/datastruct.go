package oto

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/datatypes"
)


type ValueType string
type JobStatus string

const (
	Integer   ValueType = "int"
	String    ValueType = "string"
	Tuple     ValueType = "tuple"
	FilePath  ValueType = "path"
	Float     ValueType = "float"
	IPAddress ValueType = "ip"
	Port      ValueType = "port"
	Stoped	  JobStatus = "stoped"
	Running   JobStatus = "running"
	Waiting   JobStatus = "waiting"
	Done	  JobStatus = "done"
)

type ExportedJsonConfig struct {
	Executables []Executable	`json:"executables"`
	Parameters []Parameter		`json:"parameters"`
	Commands []Command			`json:"commands"`
}

type Executable struct {
	gorm.Model
	ID        	uint   `gorm:"primaryKey;autoIncrement"`
	ExecId		string `gorm:"unique; not null" json:"exec_id" validate:"required"`
	Name        string `gorm:"not null type:string" json:"name" validate:"required"`
	Version     string `gorm:"not null type:string" json:"version" validate:"required"`
	Binary      string `gorm:"not null type:string" json:"binary" validate:"required"`
	Description string `json:"description type:string" validate:"omitempty"`
}

type Parameter struct {
	gorm.Model
	ID        		 uint   		`gorm:"primaryKey;autoIncrement"`
	Flag             string     	`gorm:"unique; not null type:string" json:"flag" validate:"required"`
	ExecId			 string			`gorm:"not null type:string" json:"exec_id" validate:"required"`
	RequiresRoot     bool       	`gorm:"not null type:bool" json:"requires_root type:bool" validate:"required"`
	RequiresValue    bool       	`gorm:"not null type:bool" json:"requires_value type:bool" validate:"required"`
	ValueType        ValueType  	`gorm:"not null" json:"value_type" validate:"omitempty"`
	DependsOn     	 datatypes.JSON `gorm:"type:json"`
	ConflictWith  	 datatypes.JSON `gorm:"type:json"`
	Description      string    		`gorm:"type:string" json:"description,omitempty"`
}

type Command struct {
	gorm.Model
	ID       		uint   			`gorm:"primaryKey;autoIncrement"`
	Name 			string			`gorm:"unique; not null type:string"`
	ExecId 			string			`gorm:"not null type:string"`
	Description 	string			`gorm:"type:string"`
	ParameterFlags  datatypes.JSON  `gorm:"type:json"`
	RequiresRoot    bool			`gorm:"not null type:bool" json:"requires_root type:bool" validate:"required"`
}

type JobCmd struct {
	Name string
	Exec *Executable
	Cmd *Command
	Results []Output
	Status JobStatus
}

type Output struct {
	Stdout string
	Stderr string
}


func NewExecutable(exec_name string, version string, binaryPath string, description string, conf *Oto) (*Executable, error) {
	execId := GetExecId(exec_name, version)
	exec := &Executable {
		ExecId: execId,
		Name: exec_name,
		Version: version,
		Binary: binaryPath,
		Description: description,
	}
	conf.Database.Create(exec)
	return exec, nil
}

func NewCommand(execId string, cmdName, description string, flags []string, conf *Oto) (*Command, error) {
	flagsBytes, err := SliceToJson(flags)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert flags to json")
	}

	var params []Parameter
	if err := conf.Database.
		Where("flag IN ?", flags).
		Find(&params).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch parameters: %w", err)
	}

	requiresRoot := false
	for _, p := range params {
		if p.RequiresRoot {
			requiresRoot = true
			break
		}
	}

	cmd := &Command{
		Name: cmdName,
		ExecId: execId,
		Description: description,
		ParameterFlags: flagsBytes,
		RequiresRoot: requiresRoot,
	}
	if err := conf.Database.Create(cmd).Error; err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}
	return cmd, nil
}

func NewParameter(execId string, flag string, description string, reqRoot bool, reqValue bool, valueType ValueType, dependsOn []string, conflictWith []string, conf *Oto) (*Parameter, error) {
	dependsOnBytes, err := SliceToJson(dependsOn)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert dependsOn to json")
	}
	conflictWithBytes, err := SliceToJson(conflictWith)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert conflictWith to json")
	}
	
	param := &Parameter{
		Flag: flag,
		ExecId: execId,
		Description: description,
		RequiresRoot: reqRoot,
		RequiresValue: reqValue,
		ValueType: valueType,
		DependsOn: dependsOnBytes,
		ConflictWith: conflictWithBytes,
	}
	conf.Database.Create(param)
	return param, nil
}
