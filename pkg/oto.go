package oto

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

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
	Workers        map[string]WorkerItem
}

type WorkerItem struct {
	WorkerID    string
	Worker      worker.Worker
	OutputError chan (error)
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

	// tmp : automigrate with gorm until we deploy atlas completly
	instance.Database.AutoMigrate(&models.Executable{}, &models.Parameter{}, &models.Command{}, &models.Job{}, &models.FlagValue{})
	return instance, nil
}

// === Add data ===

// Add and save a new executable
//
// name			: executable name
// version		: executable version
// path			: absolute path of the executable on your machine
// description	: add a few words to describe your executable
func (i *Instance) AddExecutable(name, version, path, description string) error {
	exec := models.NewExecutable(name, version, path, description)
	if err := i.Database.Save(exec).Error; err != nil {
		return fmt.Errorf("failed to save Executable: %w", err)
	}

	return nil
}


// Add and save a new parameter
//
// execName		: executable name, ( executable name - version ) i.e : "nmap - 7.94"
// flag			: parameter flag
// description	: add a few words to describe your parameter
// requiresRoot : if this parameter require to use root rigth, set to true, else to false
// requiresRoot : if this parameter require a value beside it, set to true, else to false
// valueType 	: whathever the parameter require a value or not, precise here the value he can support. Type None if this parameter can't support any values.
// requires		: add flag of parameter this one depends on.
// interfers	: add flag of parameter this one can't be mixed to.
func (i *Instance) AddParameter(ctx context.Context, execName, flag, description string, requiresRoot, requiresValue bool, valueType models.ValueType, requires, interfers []string, s *fme.Schema) error {
	// Retrieve executable
	exec, err := models.FetchExecutable(ctx, i.Database, "tag", execName)
	if err != nil {
		return err
	}

	// Retrieve dependencies
	requireToSave, err := models.FetchFlagParameters(ctx, i.Database, "flag", requires)
	if err != nil {
		return err
	}
	InterfersWithToSave, err := models.FetchFlagParameters(ctx, i.Database, "flag", interfers)
	if err != nil {
		return err
	}

	// Verify dependencies correctness
	for _, depends := range requireToSave {
		if ok, err := s.Require(flag, depends.Flag); ok == false {
			return fmt.Errorf("incorrect `requirement` dependency : %v", err)
		}
	}
	for _, Interfer := range InterfersWithToSave {
		if ok, err := s.Interfer(flag, Interfer.Flag); ok == false {
			return fmt.Errorf("incorrect `interference` dependency : %v", err)
		}
	}

	param := models.NewParameter(flag, description, exec, requiresRoot, requiresValue, valueType, InterfersWithToSave, requireToSave)
	if err := i.Database.Save(param).Error; err != nil {
		return fmt.Errorf("failed to parameter : %w", err)
	}
	return nil
}


// Add and save a new command
//
// execName		: executable tag, ( executable name - version ) i.e : "nmap - 7.94"
// flag			: parameter flag
// description	: add a few words to describe your parameter
func (i *Instance) AddCommand(ctx context.Context, execTag, cmdName, description string, root bool, flags []string, s *fme.Schema) error {
	exec, err := models.FetchExecutable(ctx, i.Database, "tag", execTag)
	if err != nil {
		return err
	}

	flagsToSave, err := models.FetchFlagParameters(ctx, i.Database, "flag", flags)
	if err != nil {
		return err
	}

	// TODO : verify root flags for parameters

	// FME : check if the given flags are valid before ingestion
	_, err = s.ValidateCombination(flags)
	if err != nil {
		return err
	}

	cmd := models.NewCommand(cmdName, description, exec, root, flagsToSave)
	if err := i.Database.Save(cmd).Error; err != nil {
		return fmt.Errorf("failed to save command: %w", err)
	}
	return nil
}


// Add and save a new job
//
// execName		: executable tag, ( executable name - version ) i.e : "nmap - 7.94"
// flag			: parameter flag
// description	: add a few words to describe your parameter
// requiresRoot : if this parameter require to use root rigth, set to true, else to false
// requiresRoot : if this parameter require a value beside it, set to true, else to false
// valueType 	: whathever the parameter require a value or not, precise here the value he can support. Type None if this parameter can't support any values.
// requires		: add flag of parameter this one depends on.
// interfers	: add flag of parameter this one can't be mixed to.
func (i *Instance) AddJob(ctx context.Context, cmdName, jobName string, flagValues map[string]string) error {
	cmd, err := models.FetchCommand(ctx, i.Database, "name", cmdName)
	if err != nil {
		return err
	}

	var flagValuesToSave []*models.FlagValue
	for flag, value := range flagValues {
		param, err := simple.GetRowBy[models.Parameter](ctx, i.Database, "flag", flag)
		if err != nil {
			return fmt.Errorf("in your flag values, you indicated a parameter %s that doesn't belong to your command %s. %v", flag, cmd.Name, err)
		}
		flagValuesToSave = append(flagValuesToSave, models.NewFlagValue(param, value))
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
		err := i.AddParameter(ctx, p.ExecutableName, p.Flag, p.Description, p.RequiresRoot, p.RequiresValue, p.ValueType, p.RequireIDs, p.InterferIDs, s)
		if err != nil {
			return fmt.Errorf("failed to add parameter %s: %w", p.Flag, err)
		}
	}
	return nil
}

func (i *Instance) AddRoutine(ctx context.Context, routineName string, jobsMap map[models.Action]string) error {
	var routine *models.Routine
	var routineMap map[models.Action]models.Job

	for action, jobName := range jobsMap {
		job, err := models.FetchJob(ctx, i.Database, "name", jobName)
		if err != nil {
			return fmt.Errorf("couldn't find job : %s. %v", job, err)
		}
		routineMap[action] = *job
	}

	routine = models.NewRoutine(routineName, routineMap, make([]models.Output, 0))

	if err := i.Database.Save(routine).Error; err != nil {
		return fmt.Errorf("failed to save job command: %w", err)
	}
	return nil
}

// Tmp function for demo only. This will be erased when Temporal we'll be fully integrated.
func (i *Instance) RunJobDemo(ctx context.Context, jobName string) (*models.Output, error) {
	job, err := models.FetchJob(ctx, i.Database, "name", jobName)
	if err != nil {
		return nil, err
	}

	var header string = job.Command.Executable.Path
	var args []string

	if job.Command.RequiresRoot {
		args = append(args, job.Command.Executable.Path)
		header = "sudo"
	}

	for _, fv := range job.FlagValues {
		// If the argument doesn't have value, we just add the flag to avoid a whitespace
		if fv.Value == "" {
			args = append(args, fv.Parameter.Flag)
		} else {
			args = append(args, fv.Parameter.Flag, fv.Value)
		}
	}

	fmt.Println(header, args)

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, header, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	return &models.Output{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}

// == FME ===
func (i *Instance) AddExecutableSchema(ctx context.Context, execName string) (*fme.Schema, error) {
	s := fme.NewSchema()

	exec, err := models.FetchExecutable(ctx, i.Database, "tag", execName)
	if err != nil {
		return nil, err
	}

	params, err := models.FetchParameters(ctx, i.Database, "id", exec.ID)
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
	i.ParamsSchema[execName] = *s
	return s, nil
}

// === Temporal ===
func (i *Instance) NewWorkerItem(workerID string) WorkerItem {
	return WorkerItem{
		WorkerID:    workerID,
		Worker:      worker.New(i.TemporalClient, workerID, worker.Options{}),
		OutputError: make(chan error),
	}
}

// Create a new worker with the given ID. The worker will then be runned concurrently.
func (i *Instance) AddWorker(workerID string) error {
	if _, ok := i.Workers[workerID]; ok {
		return fmt.Errorf("worker %s already.", workerID)
	}

	w := i.NewWorkerItem(workerID)
	i.Workers[workerID] = w

	acts := &Activities{DB: i.Database}

	w.Worker.RegisterWorkflow(WorkflowRunRoutine)
	w.Worker.RegisterActivity(acts.RunJob)

	go func() {
		if err := w.Worker.Run(worker.InterruptCh()); err != nil {
			w.OutputError <- err
		}
	}()

	return nil
}

// Stop the worker. Wait for in-going tasks to finish, and then the worker will stop.
func (i *Instance) StopWorker(workerID string) error {
	w, ok := i.Workers[workerID]
	if !ok {
		return fmt.Errorf("couldn't find workerID. First create the worker.")
	}

	w.Worker.Stop()
	return nil
}

func (i *Instance) ConfigWorker(workerID string) error {
	// TODO

	return nil
}

func (i *Instance) ConfigWorkflow(workerID string) error {
	// TODO

	return nil
}

func (i *Instance) RunWorkflow(ctx context.Context, workerID, jobName string) (*models.Output, error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("job-%s-%d", jobName, time.Now().UnixNano()),
		TaskQueue: workerID,
	}

	handle, err := i.TemporalClient.ExecuteWorkflow(ctx, workflowOptions, WorkflowRunRoutine, jobName)
	if err != nil {
		return nil, err
	}

	var result models.Output
	if err := handle.Get(ctx, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
