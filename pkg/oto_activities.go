package oto

import (
	"bytes"
	"context"
	"os/exec"

	"gorm.io/gorm"
	"github.com/Bl4omArchie/oto/models"
)

type Activities struct {
	DB *gorm.DB
}

type (
	CommandRaw struct {
		Header string
		Args []string
	}

	JobOutput struct {
		Stdout string
		Stderr string
	}
)


func (a *Activities) CraftJob(ctx context.Context, jobName string) (*CommandRaw, error) {
	job, err := models.FetchJob(ctx, a.DB, "name", jobName)
	if err != nil {
		return nil, err
	}

	header := job.Command.Binary.Path
	if job.Command.RequiresRoot {
		header = "sudo " + header
	}

	var args []string
	for _, fv := range job.FlagValues {
		args = append(args, fv.Parameter.Flag, fv.Value)
	}

	return &CommandRaw{Header: header, Args: args}, nil
}

func (a *Activities) ExecuteCommand(ctx context.Context, cmdRaw *CommandRaw) (*JobOutput, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, cmdRaw.Header, cmdRaw.Args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return &JobOutput{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}
