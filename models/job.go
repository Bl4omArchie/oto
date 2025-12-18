package models

import (
	"fmt"
	"time"
	"context"

	"gorm.io/gorm"
)


// A job is a command with pre-defined values
//
// In order to create a Job, you must specify the command name and fill structs called FlagValue which is a key-value combination between parameter flag and its value
type Job struct {
	gorm.Model
	Name       string       `gorm:"unique;not null"`
	CommandId  int          `gorm:"not null"`
	Command    *Command     `gorm:"foreignKey:CommandId"`
	FlagValues []*FlagValue `gorm:"many2many:job_flagvalues"`
}

type FlagValue struct {
	gorm.Model
	ParameterId int			`gorm:"not null;uniqueIndex:uid_flag_value"`
	Parameter	*Parameter	`gorm:"foreignKey:ParameterId"`
	Value		string		`gorm:"not null;uniqueIndex:uid_flag_value"`
}

type Output struct {
	gorm.Model
	JobId int			`gorm:"not null"`
	Job *Job			`gorm:"foreignKey:JobId"`
	Timestamp time.Time	`gorm:"not null"`
	Stdout string		`gorm:"not null"`
	Stderr string		`gorm:"not null"`
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

func NewOutput(job *Job, timestamp time.Time, stdout, stderr string) *Output {
	return &Output{
		JobId: int(job.CommandId),
		Job: job,
		Timestamp: timestamp,
		Stdout: stdout,
		Stderr: stderr,
	}
}

func FetchJob(ctx context.Context, db *gorm.DB, column, value any) (*Job, error) {
	var job Job

	err := db.WithContext(ctx).
		Preload("Command").
		Preload("Command.Executable").
		Preload("FlagValues").
		Preload("FlagValues.Parameter").
		Where(fmt.Sprintf("%s = ?", column), value).
		First(&job).Error
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func FetchJobs(ctx context.Context, db *gorm.DB, column, value any) ([]Job, error) {
	var jobs []Job

	err := db.WithContext(ctx).
		Preload("Command").
		Preload("Command.Executable").
		Preload("FlagValues").
		Preload("FlagValues.Parameter").
		Where(fmt.Sprintf("%s = ?", column), value).
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
