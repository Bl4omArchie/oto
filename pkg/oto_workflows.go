package oto

import (
	"time"

	"github.com/Bl4omArchie/oto/models"
	"go.temporal.io/sdk/workflow"
)

func WorkflowRunRoutine(ctx workflow.Context, routineName string) (*models.Output, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	routine, err := models.FetchRoutine(ctx, "")

	var output models.Output
	err = workflow.ExecuteActivity(ctx, "RunCommand", routineName).Get(ctx, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
