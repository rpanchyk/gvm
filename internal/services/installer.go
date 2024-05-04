package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rpanchyk/gvm/internal/models"
)

type Installer struct {
	Config *models.Config
}

func (i Installer) Install(version string) error {
	installDir := filepath.Join(i.Config.InstallDir, "go"+version)
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		downloader := &Downloader{Config: i.Config}
		sdk, err := downloader.Download(version)
		if err != nil {
			return fmt.Errorf("cannot download specified SDK: %w", err)
		}

		if err = i.unpackZip(sdk.FilePath, i.Config.InstallDir); err != nil {
			return fmt.Errorf("cannot unpack specified SDK: %w", err)
		}

		tempDir := filepath.Join(i.Config.InstallDir, "go")
		err = os.Rename(tempDir, installDir)
		if err != nil {
			return fmt.Errorf("cannot rename directory for specified SDK: %w", err)
		}
		fmt.Printf("SDK has been installed: %s\n", installDir)
	} else {
		fmt.Printf("SDK already installed: %s\n", installDir)
	}

	localDir := filepath.Join(i.Config.LocalDir, "go"+version)
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		for _, dir := range []string{"bin", "pkg"} {
			dirPath := filepath.Join(localDir, dir)
			if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return fmt.Errorf("cannot create dir: %s error: %w", dirPath, err)
			}
		}
		fmt.Printf("Local directory created: %s\n", localDir)
	}

	return nil
}

func (i Installer) unpackZip(src, dst string) error {
	archive, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("cannot open file: %s error: %w", src, err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s error: %w", filePath, err)
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("cannot create dir: %s error: %w", filePath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("cannot open file: %s error: %w", filePath, err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("cannot open zip file: %s error: %w", filePath, err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("cannot copy file: %s error: %w", filePath, err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
