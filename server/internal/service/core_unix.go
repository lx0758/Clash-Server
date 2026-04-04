//go:build !windows
// +build !windows

package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func getBinaryPath(workDir string) string {
	binary := ""
	files, _ := filepath.Glob(filepath.Join(workDir, "mihomo*"))
	for _, file := range files {
		if filepath.Ext(file) != ".exe" {
			binary = file
		}
	}
	if binary != "" {
		_ = os.Chmod(binary, 0755)
	}
	return binary
}

func setProcessDeathSignal(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Pdeathsig = syscall.SIGTERM
}
