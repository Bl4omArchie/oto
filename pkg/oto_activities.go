package oto

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/Bl4omArchie/oto/models"
	"gorm.io/gorm"
)

type Activities struct {
	DB *gorm.DB
}

type (
	JobOutput struct {
		Stdout string
		Stderr string
	}
)

func (a *Activities) RunJob(ctx context.Context, jobName string) (*JobOutput, error) {
	job, err := models.FetchJob(ctx, a.DB, "name", jobName)
	if err != nil {
		return nil, err
	}

	header := job.Command.Executable.Path
	if job.Command.RequiresRoot {
		header = "sudo " + header
	}

	var args []string
	for _, fv := range job.FlagValues {
		args = append(args, fv.Parameter.Flag, fv.Value)
	}

	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, header, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	return &JobOutput{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}
