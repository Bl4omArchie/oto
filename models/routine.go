package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Action string

var (
	Sequential Action = "seq"
	Concurrent Action = "conc"
)

type Routine struct {
	gorm.Model
	Name       string
	Jobs       map[Action]Job
	JobOutputs []Output
}

func NewRoutine(tag string, jobs map[Action]Job, outputs []Output) *Routine {
	return &Routine{
		Name:       tag,
		Jobs:       jobs,
		JobOutputs: outputs,
	}
}

func FetchRoutine(ctx context.Context, db *gorm.DB, column string, value any) (*Routine, error) {
	var routine Routine

	err := db.WithContext(ctx).
		Preload("Jobs").
		Preload("JobOutputs").
		Where(fmt.Sprintf("%s = ?", column), value).
		First(&routine).Error
	if err != nil {
		return nil, err
	}

	return &routine, nil
}
