package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println(os.Getwd())
	version, _, err := userChoice()
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}

	// Downloads selected version of arbelict.zip
	err = downloadBinaries("https://github.com/Dogel-ai/arbelict/releases", version)
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}
	log.Print("Binaries downloaded successfully")
}

func userChoice() (version string, packageType string, err error) {
	version = "latest"
	packageType = "both"
	scanner := bufio.NewScanner(os.Stdin)

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

func downloadBinaries(packageRoot string, version string) error {
	packageSource := fmt.Sprintf("%s/download/%s/arbelict.zip", packageRoot, version)
	if version == "latest" {
		packageSource = fmt.Sprintf("%s/latest/download/arbelict.zip", packageRoot)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
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
