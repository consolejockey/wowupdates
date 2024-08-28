package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ReleaseData struct {
	Data []struct {
		ID          int    `json:"id"`
		DateCreated string `json:"dateCreated"`
		FileName    string `json:"fileName"`
	} `json:"data"`
}

// Queries an URL and returns the response body
func queryAPI(apiURL string) ([]byte, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Retrieves the ID of the newest release from the API response
func getNewestReleaseID(jsonData []byte) (int, error) {
	var releases ReleaseData
	var newestReleaseID int
	var newestReleaseTime time.Time

	if err := json.Unmarshal(jsonData, &releases); err != nil {
		return 0, err
	}

	for _, release := range releases.Data {
		if !strings.Contains(release.FileName, "alpha") {
			releaseTime, err := time.Parse(time.RFC3339, release.DateCreated)
			if err != nil {
				return 0, err
			}

			if releaseTime.After(newestReleaseTime) {
				newestReleaseTime = releaseTime
				newestReleaseID = release.ID
			}
		}
	}
	return newestReleaseID, nil
}

// Downloads an addon given the respective download URL and a download path
func downloadFile(url, path, addon string) error {
	normalizedPath := filepath.FromSlash(path)
	err := os.MkdirAll(normalizedPath, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s.zip", normalizedPath, addon)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file for %s. Status code: %d", addon, response.StatusCode)
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	return nil
}
