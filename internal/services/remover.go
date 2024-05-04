package services

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github.com/rpanchyk/gvm/internal/models"
	"golang.org/x/sys/windows"
)

type Remover struct {
	Config *models.Config
}

func (r Remover) Remove(version string) error {
	goRootDir := filepath.Join(r.Config.InstallDir, "go"+version)
	if _, err := os.Stat(goRootDir); os.IsNotExist(err) {
		return fmt.Errorf("go%s is not installed because directory doesn't exist: %s", version, goRootDir)
	}
	goPathDir := filepath.Join(r.Config.LocalDir, "go"+version)
	if _, err := os.Stat(goPathDir); os.IsNotExist(err) {
		return fmt.Errorf("go%s is not installed because directory doesn't exist: %s", version, goPathDir)
	}

	if err := r.removeFromWindowsPathVar(goRootDir, goPathDir); err != nil {
		return fmt.Errorf("cannot update Path env var: %w", err)
	}

	if err := r.removeWindowsUserEnvVarIfEquals("GOROOT", goRootDir); err != nil {
		return fmt.Errorf("cannot remove GOROOT: %w", err)
	}
	if err := r.removeWindowsUserEnvVarIfEquals("GOPATH", goPathDir); err != nil {
		return fmt.Errorf("cannot remove GOROOT: %w", err)
	}

	// for _, dir := range []string{goRootDir, goPathDir} {
	// 	if err := os.RemoveAll(dir); err != nil {
	// 		return fmt.Errorf("cannot remove %s: %w", dir, err)
	// 	}
	// 	fmt.Printf("Directory %s has been removed\n", dir)
	// }
	return nil
}

func (r Remover) removeFromWindowsPathVar(dirs ...string) error {
	binDirs := []string{}
	for _, dir := range dirs {
		binDir := filepath.Join(dir, "bin")
		binDirs = append(binDirs, binDir)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home: %w", err)
	}

	oldPathEnvVar, ok := os.LookupEnv("Path")
	if !ok {
		return fmt.Errorf("cannot get Path env var")
	}
	fmt.Println("Old Path:", oldPathEnvVar)

	// Windows combines SystemEnv.Path and UserEnv.Path variables when we get it.
	// So, we go by splitted Path until faced with any path starting with user home directory.
	// This identifies starting of UserEnv.Path value.
	userPathEnvVar := []string{}
	for _, path := range strings.Split(oldPathEnvVar, ";") {
		if strings.HasPrefix(path, userHomeDir) {
			if slices.Contains(binDirs, path) {
				continue // actually, removing
			}
			userPathEnvVar = append(userPathEnvVar, path)
		}
	}

	normalizedUserPathEnvVar := []string{}
	for _, path := range userPathEnvVar {
		normalizedPath := strings.Replace(path, userHomeDir, "%USERPROFILE%", 1)
		normalizedUserPathEnvVar = append(normalizedUserPathEnvVar, normalizedPath)
	}

	updatedPath := strings.Join(normalizedUserPathEnvVar, ";")
	if err := r.setWindowsUserEnvVar("Path", updatedPath); err != nil {
		return fmt.Errorf("cannot update Path: %w", err)
	}

	fmt.Printf("User Path environment variable updated: %s\n", updatedPath)
	return nil
}

func (r Remover) setWindowsUserEnvVar(name, value string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("cannot get hostname: %w", err)
	}
	fmt.Printf("Hostname: %s\n", hostname)

	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user: %w", err)
	}
	fmt.Printf("Username: %s\n", user.Name)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Request Admin Permissions in Windows
	// https://gist.github.com/jerblack/d0eb182cc5a1c1d92d92a4c4fcc416c6
	// Run: setx /S <host> /U <user> GOROOT /usr/local/go
	verb := "runas"
	exe := "setx"
	args := fmt.Sprintf("/S %s /U %s %s \"%s\"", hostname, user.Name, name, value)

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	argPtr, _ := syscall.UTF16PtrFromString(args)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	showCmd := int32(0) // SW_HIDE

	if err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd); err != nil {
		return fmt.Errorf("cannot set global environment variable: %w", err)
	}

	fmt.Printf("User environment variable %s has been persisted as %s\n", name, value)
	return nil
}

func (r Remover) removeWindowsUserEnvVarIfEquals(name, value string) error {
	// fmt.Printf("Removing %s user env var\n", name)

	// envVar, ok := os.LookupEnv(name)
	_, ok := os.LookupEnv(name)
	if !ok {
		fmt.Printf("%s user env var is not found\n", name)
		return nil
	}
	// if envVar != value {
	// 	fmt.Printf("%s user env var differs from %s\n", envVar, value)
	// 	return nil
	// }
	if err := r.removeWindowsUserEnvVar(name); err != nil {
		return fmt.Errorf("cannot remove %s: %w", name, err)
	}

	fmt.Printf("%s user env var has been removed\n", name)
	return nil
}

func (r Remover) removeWindowsUserEnvVar(name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Request Admin Permissions in Windows
	// https://gist.github.com/jerblack/d0eb182cc5a1c1d92d92a4c4fcc416c6
	// Run: REG delete "HKCU\Environment" /F /V "variable_name"
	verb := "runas"
	exe := "reg"
	args := fmt.Sprintf("delete \"HKCU\\Environment\" /F /V \"%s\"", name)

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	argPtr, _ := syscall.UTF16PtrFromString(args)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	showCmd := int32(0) // SW_HIDE

	if err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd); err != nil {
		return fmt.Errorf("cannot set global environment variable: %w", err)
	}

	fmt.Printf("User environment variable %s has been removed\n", name)
	return nil
}
