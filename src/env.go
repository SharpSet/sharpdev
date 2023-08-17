package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadDotFile(dotfileURL, envName string) error {
	// Convert GitHub repo URL to the raw content URL
	rawURL := strings.Replace(dotfileURL, "github.com", "raw.githubusercontent.com", 1)
	rawURL = strings.Replace(rawURL, "/blob", "", 1)

	// Construct the URL to the sharpdev.yml in the repository
	fileURL := fmt.Sprintf("%s/main/envs/%s/sharpdev.yml", rawURL, envName)

	// Send GET request
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != 200 {
		return errors.New("Failed to fetch the file, status: " + resp.Status + ": " + fileURL)
	}

	// Ensure the destination directory exists
	destDir := "./env"
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	// Create the destination file
	destPath := filepath.Join(destDir, "sharpdev.yml")
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
