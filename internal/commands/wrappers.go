package commands

import (
	"fmt"
	"os"
	"os/exec"
)

type Command[T any] interface {
	Execute(args []string, callback func(T) error) (T, error)
}

type CommandExecCombinedOutput func([]string, func(interface{}) error) (string, error)
type CommandExec func([]string, func(interface{}) error) (string, error)
type CommandExecFunc func(string, ...string) (string, error)
type CommandExecFuncWithCallback func([]string, func(interface{}) error) (string, error)
type CommandExecFuncWithCallbackAndError func([]string, func(interface{}) error) (string, error)

func executeCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro ao executar comando: %w", err)
	}
	return string(output), nil
}
