package models

import (
	"fmt"
	"context"

	"gorm.io/gorm"
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

func FetchJob(ctx context.Context, db *gorm.DB, column, value any) (*Job, error) {
	var job Job

	err := db.WithContext(ctx).
		Preload("Command").
		Preload("FlagValues").
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
		Preload("FlagValues").
		Where(fmt.Sprintf("%s = ?", column), value).
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
