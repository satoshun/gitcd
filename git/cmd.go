package git

import (
	"os"
	"os/exec"
)

func CommitCmd(message string) (cmd *exec.Cmd) {
	args := []string{"commit", "-m", message}
	cmd = gitCmd(args)
	return
}

func InteractiveAddCmd() (cmd *exec.Cmd) {
	args := []string{"add", "-p"}
	cmd = gitCmd(args)
	return
}

func ResetCmd(hard bool) (cmd *exec.Cmd) {
	args := []string{"reset", "HEAD~"}
	if hard {
		args = append(args, "--hard")
	}

	cmd = gitCmd(args)
	return
}

func AmendCmd(message string) (cmd *exec.Cmd) {
	args := []string{"commit", "--amend", "-m", message}
	cmd = gitCmd(args)
	return
}

func gitCmd(args []string) (cmd *exec.Cmd) {
	cmd = exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return
}
