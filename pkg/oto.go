package oto

import (
	"fmt"
	"time"
	"context"

	"github.com/Bl4omArchie/fme"
	"github.com/Bl4omArchie/simple"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"gorm.io/gorm"

	"github.com/Bl4omArchie/oto/models"
)

type Instance struct {
	Database       *gorm.DB
	ParamsSchema   map[string]fme.Schema
	TemporalClient client.Client
}

type Config struct {
	PostgresDb        string `env:"POSTGRES_DB,required"`
	PostgresUser      string `env:"POSTGRES_USER,required"`
	PostgresPassword  string `env:"POSTGRES_PASSWORD,required"`
	PostgresPort      string `env:"POSTGRES_PORT,required"`
	PostgresSeed      string `env:"POSTGRES_SEED,required"`
	TemporalHost      string `env:"TEMPORAL_HOST,required"`
	TemporalNamespace string `env:"TEMPORAL_NAMESPACE,required"`
}

func NewInstanceOto(envPath string) (*Instance, error) {
	cfg, err := simple.GetEnv[Config](envPath)
	if err != nil {
		return nil, err
	}

	db, err := simple.OpenDatabase(simple.GetPostgres("localhost", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDb, cfg.PostgresPort))
	if err != nil {
		return nil, err
	}

	client, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalHost,
		Namespace: cfg.TemporalNamespace,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get the temporal client. %v", err)
	}

	instance := &Instance{
		Database:       db,
		ParamsSchema:   make(map[string]fme.Schema, 0),
		TemporalClient: client,
	}
	instance.Database.AutoMigrate(&models.Executable{}, &models.Parameter{}, &models.Command{}, &models.Job{}, &models.FlagValue{})
	return instance, nil
}

// === Add data ===

func (i *Instance) AddExecutable(name, version, executablePath, description string) error {
	bin := models.NewExecutable(name, version, executablePath, description)
	if err := i.Database.Save(bin).Error; err != nil {
		return fmt.Errorf("failed to save Executable: %w", err)
	}

	return nil
}

func (i *Instance) AddParameter(ctx context.Context, binTag, flag, description string, requiresRoot, requiresValue bool, valueType models.ValueType, Require, InterfersWith []string, s *fme.Schema) error {
	// Retrieve executable
	bin, err := models.FetchExecutable(ctx, i.Database, "tag", binTag)
	if err != nil {
		return err
	}

	// Retrieve dependencies
	RequireToSave, err := models.FetchFlagParameters(ctx, i.Database, "flag", Require)
	if err != nil {
		return err
	}
	InterfersWithToSave, err := models.FetchFlagParameters(ctx, i.Database, "flag", InterfersWith)
	if err != nil {
		return err
	}

	// Verify dependencies correctness
	for _, depends := range RequireToSave {
		if ok, err := s.Require(flag, depends.Flag); ok == false {
			return fmt.Errorf("incorrect `requirement` dependency : %v", err)
		}
	}
	for _, Interfer := range InterfersWithToSave {
		if ok, err := s.Interfer(flag, Interfer.Flag); ok == false {
			return fmt.Errorf("incorrect `interference` dependency : %v", err)
		}
	}

	param := models.NewParameter(flag, description, bin, requiresRoot, requiresValue, valueType, InterfersWithToSave, RequireToSave)
	if err := i.Database.Save(param).Error; err != nil {
		return fmt.Errorf("failed to parameter : %w", err)
	}
	return nil
}

func (i *Instance) AddCommand(ctx context.Context, binID, cmdName, description string, flags []string, s *fme.Schema) error {
	bin, err := models.FetchExecutable(ctx, i.Database, "tag", binID)
	if err != nil {
		return err
	}

	flagsToSave, err := models.FetchFlagParameters(ctx, i.Database, "flag", flags)
	if err != nil {
		return err
	}

	// FME : check if the given flags are valid before ingestion
	_, err = s.ValidateCombination(flags)
	if err != nil {
		return err
	}

	cmd := models.NewCommand(cmdName, description, bin, flagsToSave)
	if err := i.Database.Save(cmd).Error; err != nil {
		return fmt.Errorf("failed to save command: %w", err)
	}
	return nil
}

func (i *Instance) AddJob(ctx context.Context, cmdName, jobName string, flagValues map[*models.Parameter]string) error {
	// TODO : verify key-value pairs if flags are correct
	cmd, err := models.FetchCommand(ctx, i.Database, "name", cmdName)
	if err != nil {
		return err
	}

	var flagValuesToSave []*models.FlagValue
	for flag, value := range flagValues {
		flagValuesToSave = append(flagValuesToSave, models.NewFlagValue(flag, value))
	}

	job := models.NewJob(jobName, cmd, flagValuesToSave)
	if err := i.Database.Save(job).Error; err != nil {
		return fmt.Errorf("failed to save job command: %w", err)
	}
	return nil
}

func (i *Instance) ImportParameters(ctx context.Context, filename string, s *fme.Schema) error {
	params, err := simple.LoadFile[models.ParameterRaw](filename, -1, true)
	if err != nil {
		return err
	}

	for _, p := range params {
		err := i.AddParameter(
			ctx,
			p.ExecutableTag,
			p.Flag,
			p.Description,
			p.RequiresRoot,
			p.RequiresValue,
			p.ValueType,
			p.RequireIDs,
			p.InterferIDs,
			s,
		)
		if err != nil {
			return fmt.Errorf("failed to add parameter %s: %w", p.Flag, err)
		}
	}
	return nil
}

// == FME ===
func (i *Instance) AddExecutableSchema(ctx context.Context, binTag string) (*fme.Schema, error) {
	s := fme.NewSchema()

	bin, err := models.FetchExecutable(ctx, i.Database, "tag", binTag)
	if err != nil {
		return nil, err
	}

	params, err := models.FetchParameters(ctx, i.Database, "id", bin.ID)
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
	i.ParamsSchema[binTag] = *s
	return s, nil
}

// === Temporal ===
func (i *Instance) StartWorker(workerId string) error {
	w := worker.New(i.TemporalClient, workerId, worker.Options{})

	acts := &Activities{DB: i.Database}

	w.RegisterWorkflow(WorkflowRunJob)
	w.RegisterActivity(acts.CraftJob)
	w.RegisterActivity(acts.ExecuteCommand)

	return w.Run(worker.InterruptCh())
}

func (i *Instance) RunJobWorkflow(ctx context.Context, jobName string) (*JobOutput, error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("job-%s-%d", jobName, time.Now().UnixNano()),
		TaskQueue: "oto-tasks",
	}

	handle, err := i.TemporalClient.ExecuteWorkflow(ctx, workflowOptions, WorkflowRunJob, jobName)
	if err != nil {
		return nil, err
	}

	var result JobOutput
	if err := handle.Get(ctx, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
