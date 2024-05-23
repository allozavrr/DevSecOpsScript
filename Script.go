package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/mholt/archiver/v3"
)

func main() {
	gitleaksPath := os.Getenv("GITLEAKS_PATH")
	if gitleaksPath == "" {
		fmt.Println("gitleaks not found in PATH. Installing...")
		gitleaksPath = installGitleaks()
	}

	cmd := exec.Command(gitleaksPath, "--path=.", "--verbose", "--redact")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Secrets found. Commit is not allowed.")
		fmt.Println(string(output))
		os.Exit(1)
	}

	fmt.Println("No secrets found. Commit is allowed.")
}

func installGitleaks() string {
	gitleaksPath := "gitleaks"

	switch runtime.GOOS {
	case "windows":
		gitleaksPath += ".exe"
		downloadFromURL("https://github.com/gitleaks/gitleaks/releases/download/v8.8.6/gitleaks_8.8.6_windows_x64.zip", "gitleaks.zip")
		unzip("gitleaks.zip", gitleaksPath)
	case "linux":
		downloadFromURL("https://github.com/gitleaks/gitleaks/releases/download/v8.8.6/gitleaks_8.8.6_linux_x64.tar.gz", "gitleaks.tar.gz")
		untar("gitleaks.tar.gz", gitleaksPath)
	case "darwin":
		downloadFromURL("https://github.com/gitleaks/gitleaks/releases/download/v8.8.6/gitleaks_8.8.6_darwin_x64.tar.gz", "gitleaks.tar.gz")
		untar("gitleaks.tar.gz", gitleaksPath)
	default:
		fmt.Println("Unsupported operating system. Please install gitleaks manually.")
		os.Exit(1)
	}

	// Enable gitleaks for this repository
	exec.Command("git", "config", "--local", "--add", "hooks.gitleaks.enabled", "true").Run()

	return gitleaksPath
}

func downloadFromURL(url, filepath string) {
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error downloading file: status code", resp.StatusCode)
		os.Exit(1)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}

func unzip(src, dest string) {
	err := archiver.Unarchive(src, dest)
	if err != nil {
		fmt.Println("Error unzipping file:", err)
		os.Exit(1)
	}
}

func untar(src, dest string) {
	err := archiver.Unarchive(src, dest)
	if err != nil {
		fmt.Println("Error untarring file:", err)
		os.Exit(1)
	}
}
