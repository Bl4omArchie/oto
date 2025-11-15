package oto

import (
	"fmt"
	"bytes"
	"context"
	"os/exec"

	"gorm.io/gorm"
	"go.temporal.io/sdk/client"
	"github.com/Bl4omArchie/simple"
	"github.com/Bl4omArchie/oto/models"
)

type Oto struct {
	Database *gorm.DB
	TemporalClient client.Client
}

func NewInstanceOto(dbPath string) (*Oto, error) {
	db, err := simple.OpenDatabase(simple.GetSqlite(dbPath))
	if err != nil {
		return nil, err
	}

	client, err := client.Dial(client.Options{})
	if err != nil {
		return nil, err
	}

	oto := &Oto{
		Database: db,
		TemporalClient: client,
	}
	oto.Database.AutoMigrate(&models.Executable{}, &models.Parameter{}, &models.Command{}, &models.JobCommand{}, &models.FlagValue{})
	return oto, nil
}

func (oto *Oto) AddCommand(ctx context.Context, execID, cmdName, description string, flags []string) error {
	exec, err := simple.GetRowBy[models.Executable](ctx, oto.Database, "exec_id", execID)
	if err != nil {
		return fmt.Errorf("exec ID : %s, doesn't exist : %w", execID, err)
	}

	var flagsToSave []*models.Parameter
	for _, flag := range flags {
		param, err := simple.GetRowBy[models.Parameter](ctx, oto.Database, "flag", flag)
		if err != nil {
			return fmt.Errorf("param flag : %s, doesn't exist. %w", flag, err)
		}
		flagsToSave = append(flagsToSave, param)
	}

	cmd := models.NewCommand(cmdName, description, exec, flagsToSave)
	if err := oto.Database.Save(cmd).Error; err != nil {
		return fmt.Errorf("failed to save command: %w", err)
	}
	return nil
}

// TODO : verify key-value pairs if flags are correct
func (oto *Oto) AddJobCommand(ctx context.Context, execID, cmdName, jobName string, flagValues map[string]string) error {
	var header string

	exec, err := simple.GetRowBy[models.Executable](ctx, oto.Database, "exec_id", execID)
	if err != nil {
		return fmt.Errorf("exec with ID : %s, doesn't exist : %w", execID, err)
	}

	cmd, err := simple.GetRowBy[models.Command](ctx, oto.Database, "name", cmdName)
	if err != nil {
		return fmt.Errorf("command with name : %s, doesn't exist : %w", cmdName, err)
	}

	header = exec.Binary
	if cmd.RequiresRoot {
		header = "sudo " + header
	} 

	var flagValuesToSave []*models.FlagValue
	for flag, value := range flagValues {
		flagValuesToSave = append(flagValuesToSave, models.NewFlagValue(flag, value))
	}

	job := models.NewJobCommand(jobName, header, flagValuesToSave)
	if err := oto.Database.Save(job).Error; err != nil {
		return fmt.Errorf("failed to save job command: %w", err)
	}
	return nil
}

func (oto *Oto) AddExecutable(name, version, binaryPath, description string) error {
	exec := models.NewExecutable(name, version, binaryPath, description)
	if err := oto.Database.Save(exec).Error; err != nil {
		return fmt.Errorf("failed to save executable: %w", err)
	}

	return nil
}

func (oto *Oto) AddParameter(ctx context.Context, execID, flag, description string, requiresRoot, requiresValue bool, valueType models.ValueType, conflictsWith, dependsOn []string) error {
	exec, err := simple.GetRowBy[models.Executable](ctx, oto.Database, "exec_id", execID)
	if err != nil {
		return fmt.Errorf("exec with ID : %s, doesn't exist : %w", execID, err)
	}

	conflictsWithToSave, err := oto.FetchParameters(ctx, conflictsWith)
	if err != nil {
		return err
	}

	dependsOnToSave, err := oto.FetchParameters(ctx, dependsOn)
	if err != nil {
		return err
	}

	param := models.NewParameter(flag, description, exec, requiresRoot, requiresValue, valueType, conflictsWithToSave, dependsOnToSave)
	if err := oto.Database.Save(param).Error; err != nil {
		return fmt.Errorf("failed to parameter : %w", err)
	}
	return nil
}

func (oto *Oto) FetchParameters(ctx context.Context, flags []string) ([]*models.Parameter, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    var result []*models.Parameter
    for _, f := range flags {
        param, err := simple.GetRowBy[models.Parameter](ctx, oto.Database, "flag", f)
        if err != nil {
            return nil, fmt.Errorf("parameter flag %s doesn't exist: %w", f, err)
        }
        result = append(result, param)
    }
    return result, nil
}

func (oto *Oto) RunJobCommand(ctx context.Context, jobName string) (*models.RunCommandOutput, error) {
	var args []string

	job, err := simple.GetRowBy[models.JobCommand](ctx, oto.Database, "name", jobName)
	if err != nil {
		return nil, err
	}

	var flags []*models.FlagValue
	err = oto.Database.Model(&job).Association("FlagValues").Find(&flags)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve the flag values of job command : %s. %w", jobName, err)
	}

	for _, flag := range flags {
		args = append(args, flag.Flag)
		args = append(args, flag.Value)
	}

	output, err := RunCommand(ctx, job.Header, args...)
	if err != nil {
		return nil, fmt.Errorf("error while running job command : %s. %w", jobName, err)
	}
	return output, nil
}

func RunCommand(ctx context.Context, header string, args ...string) (*models.RunCommandOutput, error) {
	var stdout, stderr bytes.Buffer
	
	c := exec.CommandContext(ctx, header, args...)
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()

	return &models.RunCommandOutput{Stdout: stdout.String(), Stderr: stderr.String()}, err
}
