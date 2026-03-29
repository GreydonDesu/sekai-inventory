package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// URL and file path constants for data fetching and storage.
const (
	// CardsURL points to the latest cards data in the Sekai-World repository.
	CardsURL = "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-en-diff/refs/heads/main/cards.json"

	// CharactersURL points to the latest character data in the Sekai-World repository.
	CharactersURL = "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-en-diff/refs/heads/main/gameCharacters.json"

	// LocalCardsFile is the path where the cards data is stored locally.
	LocalCardsFile = "res/cards.json"

	// LocalCharsFile is the path where the character data is stored locally.
	LocalCharsFile = "res/gameCharacters.json"

	// MetadataFile stores information about the last data update.
	MetadataFile = "res/metadata.json"
)

// ErrNoUpdateNeeded is returned by FetchAndSaveData when the remote data
// has the same Git commit ID as the locally stored metadata and no download
// is required.
var ErrNoUpdateNeeded = errors.New("no update needed; local data is up to date")

// Metadata tracks information about fetched game data files, including
// timestamps and version information used to determine when updates
// are needed.
type Metadata struct {
	// Timestamp records when the data was last fetched, in RFC3339 format.
	Timestamp string `json:"timestamp"`

	// GitCommitID stores the commit hash of the source data repository.
	GitCommitID string `json:"gitCommitID"`

	// CardsLastUpdate stores when the cards.json file was last modified.
	CardsLastUpdate string `json:"cardsLastUpdate"`

	// CharsLastUpdate stores when the gameCharacters.json file was last modified.
	CharsLastUpdate string `json:"charsLastUpdate"`
}

// fetchFile downloads a file from the given URL and saves it to filePath.
// It returns the Last-Modified timestamp from the HTTP response, or the
// current time if the header is missing.
func fetchFile(url, filePath string) (string, error) {
	// Create an HTTP client with a timeout.
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send an HTTP GET request.
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file from URL: %w", err)
	}
	defer resp.Body.Close()

	// Check if the HTTP response status is OK.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file: received status code %d", resp.StatusCode)
	}

	// Get the Last-Modified header or use the current time.
	lastModified := resp.Header.Get("Last-Modified")
	if lastModified == "" {
		lastModified = time.Now().Format(time.RFC1123)
	}

	// Create the local file.
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy the response body to the local file.
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return lastModified, nil
}

// fetchGitCommitID retrieves the latest commit hash from the Sekai-World
// repository. It returns the SHA-1 hash of the latest commit in the main
// branch.
func fetchGitCommitID() (string, error) {
	const commitURL = "https://api.github.com/repos/Sekai-World/sekai-master-db-en-diff/commits/main"

	// Create an HTTP client with a timeout.
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send an HTTP GET request.
	resp, err := client.Get(commitURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Git commit ID: %w", err)
	}
	defer resp.Body.Close()

	// Check if the HTTP response status is OK.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch Git commit ID: received status code %d", resp.StatusCode)
	}

	// Parse the JSON response.
	var data struct {
		SHA string `json:"sha"`
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return "", fmt.Errorf("failed to parse Git commit ID: %w", err)
	}

	return data.SHA, nil
}

// LoadMetadata reads MetadataFile and returns the parsed Metadata. It returns
// an error if the file does not exist or cannot be decoded.
func LoadMetadata() (*Metadata, error) {
	file, err := os.Open(MetadataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var m Metadata
	dec := json.NewDecoder(file)
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// SaveMetadata writes a Metadata record to MetadataFile. It records the current
// time as Timestamp and stores the provided Git commit ID and Last-Modified
// values for the card and character databases.
func SaveMetadata(gitCommitID, cardsLastUpdate, charsLastUpdate string) error {
	metadata := Metadata{
		Timestamp:       time.Now().Format(time.RFC3339),
		GitCommitID:     gitCommitID,
		CardsLastUpdate: cardsLastUpdate,
		CharsLastUpdate: charsLastUpdate,
	}

	// Create or overwrite the metadata file.
	file, err := os.Create(MetadataFile)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	// Write the metadata as pretty-printed JSON.
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// ProgressCallback reports progress for long-running fetch operations.
// Stage is a human-readable description of the current step, and progress
// is a value in the range [0, 1].
type ProgressCallback func(stage string, progress float64)

// FetchAndSaveData downloads and updates all required game data files.
//
// It performs the following steps:
//
//  1. Ensure the "res" directory exists.
//  2. Fetch the latest Git commit ID from the Sekai-World repository.
//  3. Compare the commit with local metadata; if unchanged, return ErrNoUpdateNeeded.
//  4. If updated, download cards.json and gameCharacters.json.
//  5. Save metadata with updated timestamps and commit ID.
//
// If the remote Git commit ID matches the local one, FetchAndSaveData returns
// ErrNoUpdateNeeded and does not perform any downloads. When provided,
// progressCb is called with stage names and progress values between 0 and 1.
func FetchAndSaveData(progressCb ProgressCallback) error {
	// Report progress if callback is provided.
	reportProgress := func(stage string, progress float64) {
		if progressCb != nil {
			progressCb(stage, progress)
		}
	}

	// Ensure the "res" directory exists.
	if err := EnsureResDirectory(); err != nil {
		return err
	}

	// 1) Check data version via Git commit ID.
	reportProgress("Checking data version", 0.0)
	latestCommitID, err := fetchGitCommitID()
	if err != nil {
		return fmt.Errorf("error fetching Git commit ID: %v", err)
	}

	var oldMeta *Metadata
	if m, err := LoadMetadata(); err == nil {
		oldMeta = m
	}

	// If we already have metadata and the commit ID matches, skip downloads.
	if oldMeta != nil && oldMeta.GitCommitID == latestCommitID {
		reportProgress("Checking data version", 1.0)
		return ErrNoUpdateNeeded
	}
	reportProgress("Checking data version", 0.2)

	// 2) Fetch and save the cards.json file.
	reportProgress("Fetching card database", 0.2)
	cardsLastUpdate, err := fetchFile(CardsURL, LocalCardsFile)
	if err != nil {
		return fmt.Errorf("error fetching cards.json: %v", err)
	}
	reportProgress("Fetching card database", 0.5)

	// 3) Fetch and save the gameCharacters.json file.
	reportProgress("Fetching character database", 0.5)
	charsLastUpdate, err := fetchFile(CharactersURL, LocalCharsFile)
	if err != nil {
		return fmt.Errorf("error fetching gameCharacters.json: %v", err)
	}
	reportProgress("Fetching character database", 0.8)

	// 4) Save the metadata.
	reportProgress("Saving metadata", 0.8)
	if err := SaveMetadata(latestCommitID, cardsLastUpdate, charsLastUpdate); err != nil {
		return fmt.Errorf("error saving metadata: %v", err)
	}
	reportProgress("Saving metadata", 1.0)

	return nil
}
