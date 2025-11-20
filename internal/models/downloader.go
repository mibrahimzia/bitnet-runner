package models

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

type WriteCounter struct {
    Total      int64
    Downloaded int64
    OnProgress func(int64, int64) // callback
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
    n := len(p)
    wc.Downloaded += int64(n)
    if wc.OnProgress != nil {
        wc.OnProgress(wc.Downloaded, wc.Total)
    }
    return n, nil
}

// NOTE: This is the primary and only definition for DownloadStatus in this package.
// Any other definition (e.g., in manifest.go) must be removed to fix the "redeclared" compiler error.
// DownloadStatus represents the download progress and is sent through a channel.
type DownloadStatus struct {
    ModelName   string
    TotalBytes  int64
    Downloaded  int64
    Progress    float64
    IsCompleted bool
    Error       string
}

// DownloadModel fetches a GGUF file from a URL.
// progressChan is used to send updates back to the UI.
func DownloadModel(url string, filename string, progressChan chan<- DownloadStatus) (string, error) {
    // Use a simple models directory in the current working directory
    modelsDir := filepath.Join(".", "models")
    
    // Create models directory if it doesn't exist
    if err := os.MkdirAll(modelsDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create models directory: %w", err)
    }

    destPath := filepath.Join(modelsDir, filename)
    tempPath := destPath + ".tmp"

    // Create the destination file
    out, err := os.Create(tempPath)
    if err != nil {
        return "", fmt.Errorf("failed to create destination file: %w", err)
    }
    defer out.Close() // Ensure the file is closed even if errors occur

    // Get the data from the URL
    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("failed to get URL: %w", err)
    }
    defer resp.Body.Close() // Ensure the response body is closed

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("bad status: %s", resp.Status)
    }

    // Initialize progress tracking
    counter := &WriteCounter{
        Total: resp.ContentLength,
        OnProgress: func(current, total int64) {
            var progress float64 = 0
            if total > 0 {
                progress = float64(current) / float64(total) * 100
            }
            status := DownloadStatus{
                ModelName:   filename,
                TotalBytes:  total,
                Downloaded:  current,
                Progress:    progress,
                IsCompleted: false,
            }
            
            // Non-blocking send to avoid slowing down the download
            select {
            case progressChan <- status:
            default:
            }
        },
    }

    // Copy the response body to the file, while also passing it through the counter
    if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
        return "", fmt.Errorf("failed to download file: %w", err)
    }

    // Rename the temporary file to the final destination name
    if err := os.Rename(tempPath, destPath); err != nil {
        return "", fmt.Errorf("failed to rename temporary file: %w", err)
    }

    // Send a final success status
    finalStatus := DownloadStatus{
        ModelName:   filename,
        TotalBytes:  counter.Total,
        Downloaded:  counter.Downloaded,
        Progress:    100,
        IsCompleted: true,
    }
    
    // Blocking send for the final message to ensure it's received
    progressChan <- finalStatus

    return destPath, nil
}