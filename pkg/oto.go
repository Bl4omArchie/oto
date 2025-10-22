package oto

import (
	"fmt"
	"bytes"
    "os/exec"
	"encoding/json"

	"github.com/Bl4omArchie/oto/db"
)


type Oto struct {
	Database *db.OtoDB
	Jobs 	 map[string]*JobCmd
}

func OpenOto(dbPath string) (*Oto, error) {
	otoDb, err := db.OpenDatabase(dbPath)
	if err != nil {
		return nil, err
	}
	oto := &Oto{
		Database: otoDb,
		Jobs: make(map[string]*JobCmd),
	}
	oto.RefreshOto()
	return oto, nil
}

func (oto *Oto) RefreshOto() {
	oto.Database.migrate(&Executable{}, &Parameter{}, &Command{})
}

func (oto *Oto) NewJob(jobName string, executableId string, commandName string) error {
	exec, err := oto.Database.GetBy[Executable]("exec_id", executableId)
	if err != nil {
		return fmt.Errorf("Couldn't find executable : %s", executableId)
	}

	cmd, err := oto.Database.GetBy[Command]("name", commandName)
	if err != nil {
		return fmt.Errorf("Couldn't find command : %s", commandName)
	}

	job := JobCmd{
		Name: jobName,
		Exec: exec,
		Cmd: cmd,
		Status: Waiting,
	}
	oto.Jobs[jobName] = &job
	return nil
}

func (oto *Oto) ExecuteJob(jobName string, values map[string]string) error {
	var stdout, stderr bytes.Buffer
	var headerCmd string
	var cmd []string

	job := oto.Jobs[jobName]
	if job.Cmd.RequiresRoot == true {
		headerCmd = "sudo" + job.Exec.Binary
	} else {
		headerCmd = job.Exec.Binary
	}

	var flags []string
	if err := json.Unmarshal(job.Cmd.ParameterFlags, &flags); err != nil {
		return fmt.Errorf("failed to unmarshal flags: %w", err)
	}

	for _, flag := range flags {
		param, err := oto.Database.GetBy[Parameter]("flag", string(flag))
		if err != nil {
			return fmt.Errorf("Incorrect parameter flag : %w", err)
		}
		cmd = append(cmd, param.Flag)
		if val, ok := values[flag]; ok { 
			cmd = append(cmd, val) 
		}
	} 
	
	c := exec.Command(headerCmd, cmd...)
	c.Stdout = &stdout
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		return fmt.Errorf("command failed: %v\nstderr: %s", err, stderr.String())
	}

	job.Results = append(job.Results, Output{ Stdout: stdout.String(), Stderr: stderr.String() })
	return nil
}
