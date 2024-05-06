//go:build !windows

package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rpanchyk/gvm/internal/models"
)

type PlatformDefaulter struct {
	Config *models.Config
}

func (d PlatformDefaulter) Set(version string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home directory: %w", err)
	}

	filePath := filepath.Join(userHomeDir, ".profile")
	cfgString := `[ -f "$HOME/.gvm/profile" ] && . "$HOME/.gvm/profile"`
	hasGvmProfile, err := d.hasGvmProfile(filePath, cfgString)
	if err != nil {
		return fmt.Errorf("could not check gvm profile in %s: %w", filePath, err)
	}
	if !hasGvmProfile {
		if err := d.addGvmProfile(filePath, cfgString); err != nil {
			return fmt.Errorf("could not add gvm profile in %s: %w", filePath, err)
		}
		fmt.Printf("gvm profile is added to file: %s\n", filePath)
	}

	gvmProfileFile := filepath.Join(userHomeDir, ".gvm", "profile")
	if err := d.updateGvmProfile(gvmProfileFile, version); err != nil {
		return fmt.Errorf("could not update gvm profile in file %s: %w", gvmProfileFile, err)
	}

	fmt.Printf("gvm profile is updated in file %s\n", gvmProfileFile)
	return nil
}

func (d PlatformDefaulter) hasGvmProfile(filePath, cfgString string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, fmt.Errorf("could not read file %s: %w", filePath, err)
		}

		if strings.TrimSpace(line) == cfgString {
			fmt.Printf("gvm profile is found in file %s\n", filePath)
			return true, nil
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Printf("gvm profile is not found in file %s\n", filePath)
	return false, nil
}

func (d PlatformDefaulter) addGvmProfile(filePath, cfgString string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(cfgString); err != nil {
		return fmt.Errorf("could not write file %s: %w", filePath, err)
	}
	return nil
}

func (d PlatformDefaulter) updateGvmProfile(filePath, version string) error {
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory %s: %w", dirPath, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not write file %s: %w", filePath, err)
	}
	defer file.Close()

	goRoot := filepath.Join(d.Config.InstallDir, "go"+version)
	goRootEnvVar := "GOROOT=\"" + goRoot + "\""
	goPathEnvVar := "GOPATH=\"" + filepath.Join(d.Config.LocalDir, "go"+version) + "\""
	pathEnvVar := "PATH=\"$GOROOT/bin:$PATH\""

	for _, line := range []string{goRootEnvVar, goPathEnvVar, pathEnvVar} {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("could not write file %s: %w", filePath, err)
		}
	}

	return nil
}
