package actions

import (
	"fmt"
	"os/exec"
	"runtime"
)

// shellaction runs a terminal command
type ShellAction struct {
	Command string
}

func (s *ShellAction) Execute(payload interface{}) error {
	var cmd *exec.Cmd

	// detect os to run the correct shell
	if runtime.GOOS == "windows" {
		// powershell on windows
		cmd = exec.Command("powershell", "-Command", s.Command)
	} else {
		// bash on linux/mac
		cmd = exec.Command("sh", "-c", s.Command)
	}

	// run it and print output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec error: %v", err)
	}

	fmt.Printf("[action] executed: %s\n[output] %s\n", s.Command, string(output))
	return nil
}