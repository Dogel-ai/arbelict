package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Dogel-ai/arbelict/internal"
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
	packageSource := fmt.Sprintf("https://github.com/Dogel-ai/arbelict/releases/download/%s/arbelict.zip", version)
	if version == "latest" {
		packageSource = "https://github.com/Dogel-ai/arbelict/releases/latest/download/arbelict.zip"
	}
	tempDir := filepath.Join(workingDir, "temp")

	err = internal.DownloadFile(packageSource, "arbelict.zip", tempDir)
	if err != nil {
		log.Fatal("Installation failed: ", err)
	}
	log.Print("Binaries downloaded successfully.")

	log.Print("Extracting files...")
	err = internal.UnzipArchive("temp/arbelict.zip", "temp")
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

func pickFiles(packageType, workingDir string) error {
	tempDir := filepath.Join(workingDir, "temp")
	sourceDir := filepath.Join(tempDir, runtime.GOOS)

	configSource := filepath.Join(tempDir, "config.yaml")
	configDestination := filepath.Join(workingDir, "config.yaml")
	internal.CopyFile(configSource, configDestination)

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
			destinationFilename += ".exe"
		}
		destinationFile := filepath.Join(workingDir, destinationFilename)

		sourceFile := filepath.Join(sourceDir, binaryFilename)
		internal.CopyFile(sourceFile, destinationFile)
		currType = "gui"
	}

	return nil
}
