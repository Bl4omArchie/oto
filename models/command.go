package models

import (
	"context"

	"gorm.io/gorm"
	"github.com/Bl4omArchie/simple"
)

type Command struct {
	gorm.Model
	Name			string		`gorm:"unique;not null"`
	Description		string		`gorm:"type:text"`
	BinaryID		int			`gorm:"not null"`
	Binary			Binary		`gorm:"foreignKey:BinaryID"`
	RequiresRoot	bool		`gorm:"not null"`
	Parameters		[]Parameter	`gorm:"many2many:command_parameters"`
}

func FetchCommand(ctx context.Context, db *gorm.DB, field string, tag any) (*Command, error) {
	cmd, err := simple.GetRowBy[Command](ctx, db, field, tag)
	if err != nil {
		return nil, err
	}
	params, err := simple.GetTable[Parameter](ctx, db, -1)
	if err != nil {
		return nil, err
	}
	cmd.Parameters = params

	return cmd, nil
}

func NewCommand(cmdName, description string, bin *Binary, flags []Parameter) *Command {
	return &Command{
		Name: cmdName,
		BinaryID: int(bin.ID),
		Binary: *bin,
		Description: description,
		Parameters: flags,
	}
}
