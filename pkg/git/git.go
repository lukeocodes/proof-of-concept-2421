package git

import (
	"os/exec"
)

func GetFileCommit(file string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%H", "--", file)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func GetFileStage(file string) (string, error) {
	cmd := exec.Command("git", "ls-files", "--stage", file)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
