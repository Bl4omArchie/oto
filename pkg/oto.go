package oto

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/Bl4omArchie/oto/models"
	"github.com/Bl4omArchie/simple"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

type Config struct {
	Database       *gorm.DB
	TemporalClient client.Client
}

func NewInstanceOto(dbPath string) (*Config, error) {
	db, err := simple.OpenDatabase(simple.GetSqlite(dbPath))
	if err != nil {
		return nil, err
	}

	client, err := client.Dial(client.Options{})
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Database:       db,
		TemporalClient: client,
	}
	cfg.Database.AutoMigrate(&models.Binary{}, &models.Parameter{}, &models.Command{}, &models.Job{}, &models.FlagValue{})
	return cfg, nil
}

func (cfg *Config) AddCommand(ctx context.Context, binID, cmdName, description string, flags []string) error {
	bin, err := simple.GetRowBy[models.Binary](ctx, cfg.Database, "tag", binID)
	if err != nil {
		return fmt.Errorf("bin ID : %s, doesn't exist : %w", binID, err)
	}

	var flagsToSave []*models.Parameter
	for _, flag := range flags {
		param, err := simple.GetRowBy[models.Parameter](ctx, cfg.Database, "flag", flag)
		if err != nil {
			return fmt.Errorf("param flag : %s, doesn't exist. %w", flag, err)
		}
		flagsToSave = append(flagsToSave, param)
	}

	cmd := models.NewCommand(cmdName, description, bin, flagsToSave)
	if err := cfg.Database.Save(cmd).Error; err != nil {
		return fmt.Errorf("failed to save command: %w", err)
	}
	return nil
}

// TODO : verify key-value pairs if flags are correct
func (cfg *Config) AddJob(ctx context.Context, binID, cmdName, jobName string, flagValues map[string]string) error {
	var header string

	bin, err := simple.GetRowBy[models.Binary](ctx, cfg.Database, "tag", binID)
	if err != nil {
		return fmt.Errorf("binary with ID : %s, doesn't exist : %w", binID, err)
	}

	cmd, err := simple.GetRowBy[models.Command](ctx, cfg.Database, "name", cmdName)
	if err != nil {
		return fmt.Errorf("command with name : %s, doesn't exist : %w", cmdName, err)
	}

	header = bin.Path
	if cmd.RequiresRoot {
		header = "sudo " + header
	}

	var flagValuesToSave []*models.FlagValue
	for flag, value := range flagValues {
		flagValuesToSave = append(flagValuesToSave, models.NewFlagValue(flag, value))
	}

	job := models.NewJob(jobName, bin, cmd, flagValuesToSave)
	if err := cfg.Database.Save(job).Error; err != nil {
		return fmt.Errorf("failed to save job command: %w", err)
	}
	return nil
}

func (cfg *Config) AddBinary(name, version, binaryPath, description string) error {
	bin := models.NewBinary(name, version, binaryPath, description)
	if err := cfg.Database.Save(bin).Error; err != nil {
		return fmt.Errorf("failed to save Binary: %w", err)
	}

	return nil
}

func (cfg *Config) AddParameter(ctx context.Context, execID, flag, description string, requiresRoot, requiresValue bool, valueType models.ValueType, conflictsWith, dependsOn []string) error {
	exec, err := simple.GetRowBy[models.Binary](ctx, cfg.Database, "tag", execID)
	if err != nil {
		return fmt.Errorf("exec with ID : %s, doesn't exist : %w", execID, err)
	}

	conflictsWithToSave, err := cfg.FetchParameters(ctx, conflictsWith)
	if err != nil {
		return err
	}

	dependsOnToSave, err := cfg.FetchParameters(ctx, dependsOn)
	if err != nil {
		return err
	}

	param := models.NewParameter(flag, description, exec, requiresRoot, requiresValue, valueType, conflictsWithToSave, dependsOnToSave)
	if err := cfg.Database.Save(param).Error; err != nil {
		return fmt.Errorf("failed to parameter : %w", err)
	}
	return nil
}

func (cfg *Config) RunJob(ctx context.Context, jobName string) (*models.RunCommandOutput, error) {
	var header string
	var args []string

	job, err := models.GetJobFromDB(ctx, cfg.Database, jobName)
	if err != nil {
		return nil, err
	}

	var flags []*models.FlagValue
	err = cfg.Database.Model(&job).Association("FlagValues").Find(&flags)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve the flag values of job command : %s. %w", jobName, err)
	}

	if job.Command.RequiresRoot {
		header = fmt.Sprintf("sudo %s", header)
	}
	for _, flag := range flags {
		args = append(args, flag.Flag)
		args = append(args, flag.Value)
	}

	output, err := RunCommand(ctx, header, args...)
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
