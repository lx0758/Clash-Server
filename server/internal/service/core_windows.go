//go:build windows
// +build windows

package service

import (
	"os/exec"
	"path/filepath"
)

func getBinaryPath(workDir string) string {
	return filepath.Join(workDir, "mihomo.exe")
}

func setProcessDeathSignal(cmd *exec.Cmd) {
	// Windows does not support Pdeathsig
	// Process will be killed when parent exits via job object or other mechanism
}
