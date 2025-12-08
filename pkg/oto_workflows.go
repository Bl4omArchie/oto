package oto


import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func WorkflowRunJob(ctx workflow.Context, jobName string) (*JobOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var output JobOutput
	err := workflow.ExecuteActivity(ctx, "RunCommand", jobName).Get(ctx, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

