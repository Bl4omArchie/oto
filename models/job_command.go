package models

import (
	"gorm.io/gorm"
)


// TODO : fast depedency checkup for conflictWith and dependsOn

type JobCommand struct {
	gorm.Model
	Name		string			`gorm:"unique;not null" json:"name" validate:"required"`
	Header		string			`gorm:"not null" json:"header" validate:"required"`
	FlagValues	[]*FlagValue	`gorm:"many2many:jobcommand_flagvalues" json:"values" validate:"required"`
}

type FlagValue struct {
	gorm.Model
	Flag string			`gorm:"not null" json:"flag" validate:"required"`
	Value string		`gorm:"not null" json:"value" validate:"required"`
}


type RunCommandOutput struct {
	Stdout string
	Stderr string
}

func NewJobCommand(jobName, header string, flagValues []*FlagValue) *JobCommand {
	return &JobCommand{
		Name: jobName,
		Header: header,
		FlagValues: flagValues,
	}
}


func NewFlagValue(flag, value string) *FlagValue {
	return &FlagValue{
		Flag: flag,
		Value: value,
	}
}
