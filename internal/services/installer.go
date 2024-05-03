package services

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rpanchyk/gvm/internal/models"
)

type Installer struct {
	Config *models.Config
}

func (i Installer) Install(version string) error {
	listFetcher := &ListFetcher{Config: i.Config}
	sdks, err := listFetcher.Fetch()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdk, err := i.findSdk(version, sdks)
	if err != nil {
		return fmt.Errorf("cannot get specified SDK: %w", err)
	}
	fmt.Printf("Found SDK: %+v\n", *sdk)

	filePath, err := i.downloadSdk(sdk.URL, i.Config.DownloadDir)
	if err != nil {
		return fmt.Errorf("cannot download specified SDK: %w", err)
	}
	fmt.Printf("Downloaded SDK: %s\n", filePath)

	installDir := filepath.Join(i.Config.InstallDir, "go"+sdk.Version)
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		if err = i.unpackZip(filePath, i.Config.InstallDir); err != nil {
			return fmt.Errorf("cannot unpack specified SDK: %w", err)
		}

		tempDir := filepath.Join(i.Config.InstallDir, "go")
		err = os.Rename(tempDir, installDir)
		if err != nil {
			return fmt.Errorf("cannot rename specified SDK: %w", err)
		}
	} else {
		fmt.Printf("Installed SDK: %s\n", installDir)
	}

	localDir := filepath.Join(i.Config.LocalDir, "go"+sdk.Version)
	for _, dir := range []string{"bin", "pkg"} {
		dirPath := filepath.Join(localDir, dir)
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("cannot create dir: %s error: %w", dirPath, err)
		}
	}
	fmt.Printf("Local directory created: %s\n", localDir)

	return nil
}

func (i Installer) findSdk(version string, sdks []models.Sdk) (*models.Sdk, error) {
	for _, sdk := range sdks {
		if sdk.Version == version && sdk.Os == runtime.GOOS && sdk.Arch == runtime.GOARCH {
			return &sdk, nil
		}
	}
	return nil, fmt.Errorf("version %s not found", version)
}

func (i Installer) downloadSdk(fileUrl, dir string) (string, error) {
	fileName := path.Base(fileUrl)
	filePath := filepath.Join(dir, fileName)
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("File %s has been already downloaded\n", filePath)
		return filePath, nil
	}

	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("cannot get data from url: %s error: %w", fileUrl, err)
	}
	defer resp.Body.Close()

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create dir: %s error: %w", dir, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %s error: %w", filePath, err)
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("cannot save file: %s error: %w", filePath, err)
	}

	fmt.Printf("File %s has been downloaded\n", filePath)
	return filePath, nil
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
