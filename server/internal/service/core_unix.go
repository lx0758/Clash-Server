//go:build !windows
// +build !windows

package service

import (
	"os/exec"
	"path/filepath"
	"syscall"
)

func getBinaryPath(workDir string) string {
	return filepath.Join(workDir, "mihomo")
}

func setProcessDeathSignal(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Pdeathsig = syscall.SIGTERM
}
