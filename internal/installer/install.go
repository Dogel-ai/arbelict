package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}

	version, packageType, err := userChoice()
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}

	// Downloads selected version of arbelict.zip
	err = downloadBinaries("https://github.com/Dogel-ai/arbelict/releases", version, workingDir)
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}
	log.Print("Binaries downloaded successfully.")

	log.Print("Extracting files...")
	err = unzipArchive("temp/arbelict.zip", "temp")
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}
	log.Print("Files extracted successfully.")

	err = pickFiles(packageType, workingDir)
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}

	log.Print("Cleaning up...")
	err = os.RemoveAll(filepath.Join(workingDir, "temp"))
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}

	log.Print("Installation complete.")
}

func userChoice() (version string, packageType string, err error) {
	version = "latest"
	packageType = "both"
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(`



 ▄▄▄       ██▀███   ▄▄▄▄   ▓█████  ██▓     ██▓ ▄████▄  ▄▄▄█████▓
▒████▄    ▓██ ▒ ██▒▓█████▄ ▓█   ▀ ▓██▒    ▓██▒▒██▀ ▀█  ▓  ██▒ ▓▒
▒██  ▀█▄  ▓██ ░▄█ ▒▒██▒ ▄██▒███   ▒██░    ▒██▒▒▓█    ▄ ▒ ▓██░ ▒░
░██▄▄▄▄██ ▒██▀▀█▄  ▒██░█▀  ▒▓█  ▄ ▒██░    ░██░▒▓▓▄ ▄██▒░ ▓██▓ ░ 
 ▓█   ▓██▒░██▓ ▒██▒░▓█  ▀█▓░▒████▒░██████▒░██░▒ ▓███▀ ░  ▒██▒ ░ 
 ▒▒   ▓▒█░░ ▒▓ ░▒▓░░▒▓███▀▒░░ ▒░ ░░ ▒░▓  ░░▓  ░ ░▒ ▒  ░  ▒ ░░   
  ▒   ▒▒ ░  ░▒ ░ ▒░▒░▒   ░  ░ ░  ░░ ░ ▒  ░ ▒ ░  ░  ▒       ░    
  ░   ▒     ░░   ░  ░    ░    ░     ░ ░    ▒ ░░          ░      
      ░  ░   ░      ░         ░  ░    ░  ░ ░  ░ ░               
                         ░                    ░                 
						 
						 
`)

	fmt.Println("Which binaries do you want to download? (cli/gui/both):")
	fmt.Print("(default=both) ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return version, packageType, fmt.Errorf("failed reading input: %w", err)
	}
	if scanner.Text() != "" {
		packageType = scanner.Text()
	}

	fmt.Println("Version:")
	fmt.Print("(default=latest) ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return version, packageType, fmt.Errorf("failed reading input: %w", err)
	}
	if scanner.Text() != "" {
		version = scanner.Text()
	}

	return version, packageType, nil
}

func downloadBinaries(packageRoot, version, workingDir string) error {
	packageSource := fmt.Sprintf("%s/download/%s/arbelict.zip", packageRoot, version)
	if version == "latest" {
		packageSource = fmt.Sprintf("%s/latest/download/arbelict.zip", packageRoot)
	}

	tempDir := filepath.Join(workingDir, "temp")
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	filePath := filepath.Join(tempDir, "arbelict.zip")

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	log.Print("Downloading binaries...")

	resp, err := http.Get(packageSource)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received status code %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

func unzipArchive(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func pickFiles(packageType, workingDir string) error {
	tempDir := filepath.Join(workingDir, "temp")
	sourceDir := filepath.Join(tempDir, runtime.GOOS)

	configSource := filepath.Join(tempDir, "config.yaml")
	configDestination := filepath.Join(workingDir, "config.yaml")
	copyFile(configSource, configDestination)

	currType := packageType
	i := 1

	if packageType == "both" {
		currType = "cli"
		i = 0
	}

	for ; i < 2; i++ {
		binaryFilename := fmt.Sprintf("%s-%s-%s", currType, runtime.GOOS, runtime.GOARCH)
		destinationFilename := fmt.Sprintf("arbelict-%s", currType)

		if runtime.GOOS == "windows" {
			binaryFilename += ".exe"
		}
		destinationFile := filepath.Join(workingDir, destinationFilename)

		sourceFile := filepath.Join(sourceDir, binaryFilename)
		copyFile(sourceFile, destinationFile)
		currType = "gui"
	}

	return nil
}

func copyFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed opening file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed copying file: %w", err)
	}

	err = destFile.Chmod(0755)
	if err != nil {
		return fmt.Errorf("failed setting file permissions: %w", err)
	}

	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed syncing file: %w", err)
	}

	return nil
}
