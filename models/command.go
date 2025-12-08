package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Command struct {
	gorm.Model
	Name         string      `gorm:"unique;not null"`
	Description  string      `gorm:"type:text"`
	ExecutableID int         `gorm:"not null"`
	Executable   Executable  `gorm:"foreignKey:ExecutableID"`
	RequiresRoot bool        `gorm:"not null"`
	Parameters   []Parameter `gorm:"many2many:command_parameters"`
}

func NewCommand(cmdName, description string, exec *Executable, flags []Parameter) *Command {
	return &Command{
		Name:         cmdName,
		ExecutableID: int(exec.ID),
		Executable:   *exec,
		Description:  description,
		Parameters:   flags,
	}
}

// FetchCommand returns the first command corresponding to the given column and value.
func FetchCommand(ctx context.Context, db *gorm.DB, column string, value any) (*Command, error) {
	var cmd Command

	err := db.WithContext(ctx).
		Preload("Parameters").
		Where(fmt.Sprintf("%s = ?", column), value).
		First(&cmd).Error
	if err != nil {
		return nil, err
	}

	return &cmd, nil
}

// FetchCommands returns every commands corresponding to the given column and value.
func FetchCommands(ctx context.Context, db *gorm.DB, column string, value any) ([]Command, error) {
	var cmds []Command

	err := db.WithContext(ctx).
		Preload("Parameter").
		Where(fmt.Sprintf("%s = ?", column), value).
		Find(&cmds).Error
	if err != nil {
		return nil, err
	}

	return cmds, nil
}
