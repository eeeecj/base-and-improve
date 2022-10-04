package scheduler

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteTask(execCmd string) {
	execParts := strings.SplitN(execCmd, " ", 2)

	execName := execParts[0]
	execParams := ""
	fmt.Println(execParts)
	if len(execParts) > 1 {
		execParams = execParts[1]
	}

	cmd := exec.Command(execName, execParams)

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
