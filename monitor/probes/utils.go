package probes

import (
	"bytes"
	"os/exec"

	log "github.com/cihub/seelog"
)

// CommandResult stores the result of an executed command
type CommandResult struct {
	Stdout     string
	Stderr     string
	StatusCode int
}

type commandRunner interface {
	runCommand(command []string) CommandResult
}

// LinuxCommandRunner executes the given command on a Linux OS
type LinuxCommandRunner struct{}

func (runner LinuxCommandRunner) runCommand(command []string) CommandResult {
	commandName := command[0]
	commandArgs := command[1:]
	log.Debugf("Running command %q with arguments %q", commandName, commandArgs)
	cmd := exec.Command(commandName, commandArgs...)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()
	return CommandResult{
		Stdout:     out.String(),
		Stderr:     err.String(),
		StatusCode: cmd.ProcessState.ExitCode(),
	}
}
