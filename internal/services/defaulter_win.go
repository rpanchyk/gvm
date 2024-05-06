//go:build windows

package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rpanchyk/gvm/internal/models"
)

type PlatformDefaulter struct {
	Config *models.Config
}

func (d PlatformDefaulter) Set(version string) error {
	goRootDir := filepath.Join(d.Config.InstallDir, "go"+version)
	if _, err := os.Stat(goRootDir); os.IsNotExist(err) {
		return fmt.Errorf("go%s is not installed because directory doesn't exist: %s", version, goRootDir)
	}
	goPathDir := filepath.Join(d.Config.LocalDir, "go"+version)
	if _, err := os.Stat(goPathDir); os.IsNotExist(err) {
		return fmt.Errorf("go%s is not installed because directory doesn't exist: %s", version, goPathDir)
	}

	if err := d.setUserEnvVar("GOROOT", goRootDir); err != nil {
		return fmt.Errorf("cannot set GOROOT: %w", err)
	}
	if err := d.setUserEnvVar("GOPATH", goPathDir); err != nil {
		return fmt.Errorf("cannot set GOPATH: %w", err)
	}

	goBinDirs := []string{"%GOROOT%\\bin", "%GOPATH%\\bin"}
	if err := d.updatePathUserEnvVar(goBinDirs); err != nil {
		return fmt.Errorf("cannot add %s to Path: %w", goBinDirs, err)
	}

	return nil
}

func (d PlatformDefaulter) setUserEnvVar(name, value string) error {
	command := fmt.Sprintf("[Environment]::SetEnvironmentVariable(\"%s\",\"%s\",\"User\")", name, value)
	if _, err := d.runPowershellCommand(command); err != nil {
		return fmt.Errorf("powershell error: %w", err)
	}

	fmt.Printf("Set user environment variable: %s=%s\n", name, value)
	return nil
}

func (d PlatformDefaulter) getUserEnvVar(name string) (string, error) {
	command := fmt.Sprintf("[System.Environment]::GetEnvironmentVariable(\"%s\",\"User\")", name)
	value, err := d.runPowershellCommand(command)
	if err != nil {
		return "", fmt.Errorf("powershell error: %w", err)
	}
	fmt.Printf("Get user environment variable: %s=%s\n", name, value)
	return value, nil
}

func (d PlatformDefaulter) runPowershellCommand(command string) (string, error) {
	powershell, err := exec.LookPath("powershell.exe")
	if err != nil {
		return "", fmt.Errorf("cannot find powershell in Path: %w", err)
	}

	args := []string{"-NoProfile", "-NonInteractive", command}
	cmd := exec.Command(powershell, args...)

	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("executing powershell command error: %w", err)
	}
	stdOut, stdErr := stdOutBuf.String(), stdErrBuf.String()
	if stdErr != "" {
		return "", fmt.Errorf("executing powershell command failed: %s", stdErr)
	}
	fmt.Printf("Powershell command output: %s", stdOut)
	return stdOut, nil
}

func (d PlatformDefaulter) updatePathUserEnvVar(values []string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home: %w", err)
	}

	oldPathEnvVar, err := d.getUserEnvVar("Path")
	if err != nil {
		return fmt.Errorf("powershell error: %w", err)
	}

	pathEnvVar := []string{}
	for _, path := range strings.Split(oldPathEnvVar, ";") {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if strings.HasPrefix(path, d.Config.InstallDir) || strings.HasPrefix(path, d.Config.LocalDir) {
			continue
		}
		if strings.HasPrefix(path, "%GOROOT%") || strings.HasPrefix(path, "%GOPATH%") {
			continue
		}
		pathEnvVar = append(pathEnvVar, path)
	}
	pathEnvVar = append(pathEnvVar, values...)

	normalizedPathEnvVar := []string{}
	for _, path := range pathEnvVar {
		normalizedPath := strings.Replace(path, userHomeDir, "%USERPROFILE%", 1)
		normalizedPathEnvVar = append(normalizedPathEnvVar, normalizedPath)
	}

	updatedPath := strings.Join(normalizedPathEnvVar, ";")
	if err := d.setUserEnvVar("Path", updatedPath); err != nil {
		return fmt.Errorf("cannot update Path: %w", err)
	}

	fmt.Printf("User Path environment variable updated: %s\n", updatedPath)
	return nil
}
