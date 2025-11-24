package models

import (
	"context"

	"gorm.io/gorm"
	"github.com/Bl4omArchie/simple"
)


type Job struct {
	gorm.Model
	Name       string       `gorm:"unique;not null"`
	CommandId  int          `gorm:"not null"`
	Command    *Command     `gorm:"foreignKey:CommandId"`
	FlagValues []*FlagValue `gorm:"many2many:job_flagvalues"`
}

type FlagValue struct {
	gorm.Model
	ParameterId  int		`gorm:"not null"`
	Parameter    *Parameter	`gorm:"foreignKey:ParameterId"`
	Value string			`gorm:"not null"`
}

type JobOutput struct {
	Stdout string
	Stderr string
}


func NewJob(jobName string, cmd *Command, flagValues []*FlagValue) *Job {
	return &Job{
		Name: jobName,
		CommandId: int(cmd.ID),
		Command: cmd,
		FlagValues: flagValues,
	}
}

func NewFlagValue(param *Parameter, value string) *FlagValue {
	return &FlagValue{
		ParameterId: int(param.ID),
		Parameter: param,
		Value: value,
	}
}

func FetchJob(ctx context.Context, db *gorm.DB, column, jobName string) (*Job, error) {
	job, err := simple.GetRowBy[Job](ctx, db, column, jobName)
	if err != nil {
		return nil, err
	}

	cmd, err := FetchCommand(ctx, db, "ID", job.CommandId)
	if err != nil {
		return nil, err
	}
	job.Command = cmd

	return job, nil
}
