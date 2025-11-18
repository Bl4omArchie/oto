package models

import (
	"context"

	"gorm.io/gorm"
	"github.com/Bl4omArchie/simple"
)


type Job struct {
	gorm.Model
	Name       string       `gorm:"unique;not null"`
	BinaryId   int          `gorm:"not null"`
	Binary     *Binary      `gorm:"foreignKey:BinaryId"`
	CommandId  int          `gorm:"not null"`
	Command    *Command     `gorm:"foreignKey:CommandId"`
	FlagValues []*FlagValue `gorm:"many2many:job_flagvalues"`
}

type FlagValue struct {
	gorm.Model
	Flag string			`gorm:"not null"`
	Value string		`gorm:"not null"`
}

type RunCommandOutput struct {
	Stdout string
	Stderr string
}

func NewJob(jobName string, bin *Binary, cmd *Command, flagValues []*FlagValue) *Job {
	return &Job{
		Name: jobName,
		BinaryId: int(bin.ID),
		Binary: bin,
		CommandId: int(cmd.ID),
		Command: cmd,
		FlagValues: flagValues,
	}
}

func FetchJob(ctx context.Context, db *gorm.DB, jobName string) (*Job, error) {
	job, err := simple.GetRowBy[Job](ctx, db, "name", jobName)
	if err != nil {
		return nil, err
	}

	bin, err := FetchBinary(ctx, db, "ID", job.BinaryId)
	if err != nil {
		return nil, err
	}
	job.Binary = bin

	cmd, err := FetchCommand(ctx, db, "ID", job.CommandId)
	if err != nil {
		return nil, err
	}
	job.Command = cmd

	return job, nil
}

func NewFlagValue(flag, value string) *FlagValue {
	return &FlagValue{
		Flag: flag,
		Value: value,
	}
}
