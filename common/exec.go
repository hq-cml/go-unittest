package common

import (
	"fmt"
	"errors"
	"os/exec"
)

func Exec(cmd string, args ...string) (string, error) {
	cmdPath, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Errorf("exec.LookPath err: %v, cmd: %s", err, cmd)
		return "", errors.New("any")
	}

	var output []byte
	output, err = exec.Command(cmdPath, args...).CombinedOutput()
	if err != nil {
		fmt.Errorf("exec.Command.CombinedOutput err: %v, cmd: %s", err, cmd)
		return "", errors.New("any")
	}
	fmt.Println("CMD[", cmdPath, "]ARGS[", args, "]OUT[", string(output), "]")
	return string(output), nil
}
