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

type Instance struct {
	Database       *gorm.DB
	ParamsSchema   map[string]Schema
	TemporalClient client.Client
}

func NewInstanceOto() (*Instance, error) {
	cfg, err := LoadOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	db, err := simple.OpenDatabase("postgresql", cfg.PostgresqlDsn)
	if err != nil {
		return nil, err
	}

	client, err := client.Dial(client.Options{})
	if err != nil {
		return nil, err
	}

	instance := &Instance{
		Database:       db,
		ParamsSchema:   make(map[string]Schema, 0),
		TemporalClient: client,
	}
	instance.Database.AutoMigrate(&models.Binary{}, &models.Parameter{}, &models.Command{}, &models.Job{}, &models.FlagValue{})
	return instance, nil
}

func (cfg *Instance) AddCommand(ctx context.Context, binID, cmdName, description string, flags []string, s *Schema) error {
	bin, err := models.FetchBinary(ctx, cfg.Database, "tag", binID)
	if err != nil {
		return err
	}

	flagsToSave, err := models.FetchFlagParameters(ctx, cfg.Database, "flag", flags)
	if err != nil {
		return err
	}

	// FME : check if the given flags are valid before ingestion
	_, err = s.ValidateCombination(flags)
	if err != nil {
		return err
	}

	cmd := models.NewCommand(cmdName, description, bin, flagsToSave)
	if err := cfg.Database.Save(cmd).Error; err != nil {
		return fmt.Errorf("failed to save command: %w", err)
	}
	return nil
}

// TODO : verify key-value pairs if flags are correct
func (cfg *Instance) AddJob(ctx context.Context, cmdName, jobName string, flagValues map[*models.Parameter]string) error {
	cmd, err := models.FetchCommand(ctx, cfg.Database, "name", cmdName)
	if err != nil {
		return err
	}

	var flagValuesToSave []*models.FlagValue
	for flag, value := range flagValues {
		flagValuesToSave = append(flagValuesToSave, models.NewFlagValue(flag, value))
	}

	job := models.NewJob(jobName, cmd, flagValuesToSave)
	if err := cfg.Database.Save(job).Error; err != nil {
		return fmt.Errorf("failed to save job command: %w", err)
	}
	return nil
}

func (cfg *Instance) AddBinary(name, version, binaryPath, description string) error {
	bin := models.NewBinary(name, version, binaryPath, description)
	if err := cfg.Database.Save(bin).Error; err != nil {
		return fmt.Errorf("failed to save Binary: %w", err)
	}

	return nil
}

func (cfg *Instance) AddParameter(ctx context.Context, binTag, flag, description string, requiresRoot, requiresValue bool, valueType models.ValueType, Require, InterfersWith []string, s *Schema) error {
	bin, err := models.FetchBinary(ctx, cfg.Database, "tag", binTag)
	if err != nil {
		return err
	}

	RequireToSave, err := models.FetchFlagParameters(ctx, cfg.Database, "flag", Require)
	if err != nil {
		return err
	}

	for _, depends := range RequireToSave {
		s.Require(flag, depends.Flag)
	}

	InterfersWithToSave, err := models.FetchFlagParameters(ctx, cfg.Database, "flag", InterfersWith)
	if err != nil {
		return err
	}

	for _, Interfer := range InterfersWithToSave {
		s.Interfer(flag, Interfer.Flag)
	}

	err = s.ValidateSchema()
	if err != nil {
		return err
	}

	param := models.NewParameter(flag, description, bin, requiresRoot, requiresValue, valueType, InterfersWithToSave, RequireToSave)
	if err := cfg.Database.Save(param).Error; err != nil {
		return fmt.Errorf("failed to parameter : %w", err)
	}
	return nil
}

func (cfg *Instance) AddBinarySchema(ctx context.Context, binTag string) (*Schema, error) {
	s := NewSchema()

	bin, err := models.FetchBinary(ctx, cfg.Database, "tag", binTag)
	if err != nil {
		return nil, err
	}

	params, err := models.FetchParameters(ctx, cfg.Database, "binary_id", bin.ID)
	if err != nil {
		return nil, err
	}

	for _, param := range params {
		for _, dependency := range param.Require {
			s.Require(param.Flag, dependency.Flag)
		}

		for _, Interfer := range param.Interfer {
			s.Interfer(param.Flag, Interfer.Flag)
		}
	}

	err = s.ValidateSchema()
	if err != nil {
		return nil, err
	}
	cfg.ParamsSchema[binTag] = *s
	return s, nil
}

func (cfg *Instance) RunJob(ctx context.Context, jobName string) (*models.JobOutput, error) {
	var header string
	var args []string

	job, err := models.FetchJob(ctx, cfg.Database, "name", jobName)
	if err != nil {
		return nil, err
	}

	var flags []*models.FlagValue
	err = cfg.Database.Model(&job).Association("FlagValues").Find(&flags)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve the flag values of job command : %s. %w", jobName, err)
	}

	header = job.Command.Binary.Path
	if job.Command.RequiresRoot {
		header = fmt.Sprintf("sudo %s", header)
	}
	for _, flag := range flags {
		args = append(args, flag.Parameter.Flag)
		args = append(args, flag.Value)
	}

	output, err := RunCommand(ctx, header, args...)
	if err != nil {
		return nil, fmt.Errorf("error while running job command : %s. %w", jobName, err)
	}
	return output, nil
}

func RunCommand(ctx context.Context, header string, args ...string) (*models.JobOutput, error) {
	var stdout, stderr bytes.Buffer

	c := exec.CommandContext(ctx, header, args...)
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()

	return &models.JobOutput{Stdout: stdout.String(), Stderr: stderr.String()}, err
}
