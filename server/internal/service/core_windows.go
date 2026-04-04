//go:build windows
// +build windows

package service

import (
	"os/exec"
	"path/filepath"
)

func getBinaryPath(workDir string) string {
	binary := ""
	files, _ := filepath.Glob(filepath.Join(workDir, "mihomo*"))
	for _, file := range files {
		if filepath.Ext(file) == ".exe" {
			binary = file
		}
	}
	return binary
}

func setProcessDeathSignal(cmd *exec.Cmd) {
	// Windows does not support Pdeathsig
	// Process will be killed when parent exits via job object or other mechanism
}
