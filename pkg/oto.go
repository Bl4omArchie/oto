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
	Database		*gorm.DB
	ParamsSchema	map[string]Schema
	TemporalClient 	client.Client
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
		ParamsSchema:	make(map[string]Schema, 0),
		TemporalClient: client,
	}
	cfg.Database.AutoMigrate(&models.Binary{}, &models.Parameter{}, &models.Command{}, &models.Job{}, &models.FlagValue{})
	return cfg, nil
}

func (cfg *Config) AddCommand(ctx context.Context, binID, cmdName, description string, flags []string, s *Schema) error {
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
func (cfg *Config) AddJob(ctx context.Context, binID, cmdName, jobName string, flagValues map[string]string) error {
	var header string

	bin, err := models.FetchBinary(ctx, cfg.Database, "tag", binID)
	if err != nil {
		return err
	}

	cmd, err := models.FetchCommand(ctx, cfg.Database, "name", cmdName)
	if err != nil {
		return err
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

func (cfg *Config) AddParameter(ctx context.Context, binTag, flag, description string, requiresRoot, requiresValue bool, valueType models.ValueType, Require, InterfersWith []string, s *Schema) error {
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

func (cfg *Config) AddBinarySchema(ctx context.Context, binTag string) (*Schema, error) {
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

func (cfg *Config) RunJob(ctx context.Context, jobName string) (*models.RunCommandOutput, error) {
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

	header = job.Binary.Path
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
