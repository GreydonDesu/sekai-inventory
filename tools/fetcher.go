package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// URL and file path constants for data fetching and storage
const (
	// CardsURL points to the latest cards data in the Sekai-World repository
	CardsURL = "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-en-diff/refs/heads/main/cards.json"
	// CharactersURL points to the latest character data in the Sekai-World repository
	CharactersURL = "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-en-diff/refs/heads/main/gameCharacters.json"
	// LocalCardsFile is the path where the cards data is stored locally
	LocalCardsFile = "res/cards.json"
	// LocalCharsFile is the path where the character data is stored locally
	LocalCharsFile = "res/gameCharacters.json"
	// MetadataFile stores information about the last data update
	MetadataFile = "res/metadata.json"
)

// Metadata tracks information about the fetched game data files including
// timestamps and version information to help determine when updates are needed.
type Metadata struct {
	// Timestamp records when the data was last fetched
	Timestamp string `json:"timestamp"`
	// GitCommitID stores the commit hash of the source data repository
	GitCommitID string `json:"gitCommitID"`
	// CardsLastUpdate stores when the cards.json file was last modified
	CardsLastUpdate string `json:"cardsLastUpdate"`
	// CharsLastUpdate stores when the gameCharacters.json file was last modified
	CharsLastUpdate string `json:"charsLastUpdate"`
}

// fetchFile downloads a file from a given URL and saves it locally.
// It handles HTTP requests with appropriate timeouts and error handling.
//
// Parameters:
//   - url: The source URL to download from
//   - filePath: The local path where the file should be saved
//
// Returns:
//   - The Last-Modified timestamp from the HTTP response
//   - An error if the download or save operations fail
func fetchFile(url, filePath string) (string, error) {
	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send an HTTP GET request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file from URL: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check if the HTTP response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file: received status code %d", resp.StatusCode)
	}

	// Get the Last-Modified header or use current time
	lastModified := resp.Header.Get("Last-Modified")
	if lastModified == "" {
		lastModified = time.Now().Format(time.RFC1123)
	}

	// Create the local file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close() // Ensure the file is closed

	// Copy the response body to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return lastModified, nil
}

// fetchGitCommitID retrieves the latest commit hash from the Sekai-World repository.
// This helps track which version of the game data is currently being used.
//
// Returns:
//   - The SHA-1 hash of the latest commit in the main branch
//   - An error if the GitHub API request fails
func fetchGitCommitID() (string, error) {
	const commitURL = "https://api.github.com/repos/Sekai-World/sekai-master-db-en-diff/commits/main"

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send an HTTP GET request
	resp, err := client.Get(commitURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Git commit ID: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check if the HTTP response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch Git commit ID: received status code %d", resp.StatusCode)
	}

	// Parse the JSON response
	var data struct {
		SHA string `json:"sha"`
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return "", fmt.Errorf("failed to parse Git commit ID: %w", err)
	}

	return data.SHA, nil
}

// SaveMetadata stores update information in the metadata file.
// This includes timestamps and version information that can be used to
// determine when the local data needs to be updated.
//
// Parameters:
//   - gitCommitID: The current commit hash from the data repository
//   - cardsLastUpdate: Last-Modified timestamp for the cards data
//   - charsLastUpdate: Last-Modified timestamp for the character data
//
// Returns an error if writing the metadata file fails
func SaveMetadata(gitCommitID, cardsLastUpdate, charsLastUpdate string) error {
	metadata := Metadata{
		Timestamp:       time.Now().Format(time.RFC3339),
		GitCommitID:     gitCommitID,
		CardsLastUpdate: cardsLastUpdate,
		CharsLastUpdate: charsLastUpdate,
	}

	// Create or overwrite the metadata file
	file, err := os.Create(MetadataFile)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close() // Ensure the file is closed

	// Write the metadata as JSON
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// ProgressCallback is a function that receives progress updates during data fetching
type ProgressCallback func(stage string, progress float64)

// FetchAndSaveData downloads and updates all required game data files.
// This function performs the following steps:
//  1. Creates the resources directory if it doesn't exist
//  2. Downloads the latest cards.json and gameCharacters.json files
//  3. Retrieves the current Git commit ID for version tracking
//  4. Updates the metadata file with current timestamps
//
// Parameters:
//   - progressCb: Optional callback function to report progress (can be nil)
//
// Returns an error if any step in the update process fails
func FetchAndSaveData(progressCb ProgressCallback) error {
	// Report progress if callback is provided
	reportProgress := func(stage string, progress float64) {
		if progressCb != nil {
			progressCb(stage, progress)
		}
	}

	// Ensure the "res" directory exists
	if err := EnsureResDirectory(); err != nil {
		return err
	}

	// Fetch and save the cards.json file
	reportProgress("Fetching card database", 0.0)
	cardsLastUpdate, err := fetchFile(CardsURL, LocalCardsFile)
	if err != nil {
		return fmt.Errorf("error fetching cards.json: %v", err)
	}
	reportProgress("Fetching card database", 0.35)

	// Fetch and save the gameCharacters.json file
	reportProgress("Fetching character database", 0.35)
	charsLastUpdate, err := fetchFile(CharactersURL, LocalCharsFile)
	if err != nil {
		return fmt.Errorf("error fetching gameCharacters.json: %v", err)
	}
	reportProgress("Fetching character database", 0.70)

	// Fetch the latest Git commit ID
	reportProgress("Checking data version", 0.70)
	gitCommitID, err := fetchGitCommitID()
	if err != nil {
		return fmt.Errorf("error fetching Git commit ID: %v", err)
	}
	reportProgress("Checking data version", 0.85)

	// Save the metadata
	reportProgress("Saving metadata", 0.85)
	if err := SaveMetadata(gitCommitID, cardsLastUpdate, charsLastUpdate); err != nil {
		return fmt.Errorf("error saving metadata: %v", err)
	}
	reportProgress("Saving metadata", 1.0)

	return nil
}
