package models

import (
	"gorm.io/gorm"
)

// Major work todo : fast depedency checkup for conflictWith and dependsOn

type Command struct {
	gorm.Model
	Name           string         `gorm:"unique;not null" json:"name" validate:"required"`
	Description    string         `gorm:"type:text" json:"description" validate:"omitempty"`
	ExecID         int         	  `gorm:"not null" json:"exec_id" validate:"required"`
	Exec		   Executable
	RequiresRoot   bool           `gorm:"not null" json:"requires_root" validate:"required"`
	ParameterFlags []*Parameter   `gorm:"many2many:command_parameters" json:"parameter_flags" validate:"required"`
}

func NewCommand(cmdName, description string, exec *Executable, flags []*Parameter) *Command {
	return &Command{
		Name: cmdName,
		ExecID: int(exec.ID),
		Exec: *exec,
		Description: description,
		ParameterFlags: flags,
	}
}
