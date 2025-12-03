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

	var cmdRaw CommandRaw
	err := workflow.ExecuteActivity(ctx, "CraftJob", jobName).Get(ctx, &cmdRaw)
	if err != nil {
		return nil, err
	}

	var output JobOutput
	err = workflow.ExecuteActivity(ctx, "ExecuteCommand", cmdRaw.Header, cmdRaw.Args).Get(ctx, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
