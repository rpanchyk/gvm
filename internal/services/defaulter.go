package services

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"

	"github.com/rpanchyk/gvm/internal/models"
	"golang.org/x/sys/windows"
)

type Defaulter struct {
	Config *models.Config
}

func (d Defaulter) Set(version string) error {
	if runtime.GOOS == "windows" {
		return d.setWindows(version)
	}
	return d.setUnix(version)
}

// How to edit, clear, and delete environment variables in Windows
// https://www.digitalcitizen.life/remove-edit-clear-environment-variables/
func (d Defaulter) setWindows(version string) error {
	goRootDir := filepath.Join(d.Config.InstallDir, "go"+version)
	if _, err := os.Stat(goRootDir); os.IsNotExist(err) {
		return fmt.Errorf("go%s is not installed because directory doesn't exist: %s", version, goRootDir)
	}
	if err := d.setWindowsUserEnvVar("GOROOT", goRootDir); err != nil {
		return fmt.Errorf("cannot set GOROOT: %w", err)
	}

	goPathDir := filepath.Join(d.Config.LocalDir, "go"+version)
	if _, err := os.Stat(goPathDir); os.IsNotExist(err) {
		return fmt.Errorf("go%s is not installed because directory doesn't exist: %s", version, goPathDir)
	}
	if err := d.setWindowsUserEnvVar("GOPATH", goPathDir); err != nil {
		return fmt.Errorf("cannot set GOPATH: %w", err)
	}

	goBinDirs := []string{"%GOROOT%\\bin", "%GOPATH%\\bin"}
	if err := d.updateWindowsPathVar(goBinDirs); err != nil {
		return fmt.Errorf("cannot add %s to Path: %w", goBinDirs, err)
	}

	return nil
}

func (d Defaulter) setWindowsUserEnvVar(name, value string) error {
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

func (d Defaulter) updateWindowsPathVar(values []string) error {
	listFetcher := &ListFetcher{Config: d.Config}
	sdks, err := listFetcher.FetchAll()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdkBinDirs := []string{}
	for _, sdk := range sdks {
		installBinDir := filepath.Join(d.Config.InstallDir, "go"+sdk.Version, "bin")
		localBinDir := filepath.Join(d.Config.LocalDir, "go"+sdk.Version, "bin")
		sdkBinDirs = append(sdkBinDirs, installBinDir, localBinDir)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home: %w", err)
	}

	oldPathEnvVar, ok := os.LookupEnv("Path")
	if !ok {
		return fmt.Errorf("cannot get Path env var")
	}
	fmt.Println("Path:", oldPathEnvVar)

	// Windows combines SystemEnv.Path and UserEnv.Path variables when we get it.
	// So, we go by splitted Path until faced with any path starting with user home directory.
	// This identifies starting of UserEnv.Path value.
	userPathEnvVar := []string{}
	for _, path := range strings.Split(oldPathEnvVar, ";") {
		if strings.HasPrefix(path, userHomeDir) {
			if slices.Contains(sdkBinDirs, path) {
				continue // actually, removing
			}
			userPathEnvVar = append(userPathEnvVar, path)
		}
	}
	userPathEnvVar = append(userPathEnvVar, values...)

	normalizedUserPathEnvVar := []string{}
	for _, path := range userPathEnvVar {
		normalizedPath := strings.Replace(path, userHomeDir, "%USERPROFILE%", 1)
		normalizedUserPathEnvVar = append(normalizedUserPathEnvVar, normalizedPath)
	}

	updatedPath := strings.Join(normalizedUserPathEnvVar, ";")
	if err := d.setWindowsUserEnvVar("Path", updatedPath); err != nil {
		return fmt.Errorf("cannot update Path: %w", err)
	}

	fmt.Printf("User Path environment variable updated: %s\n", updatedPath)
	return nil
}

func (d Defaulter) setUnix(version string) error {
	panic("not implemented")
}
