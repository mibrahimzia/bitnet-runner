package models

import "time"

// ModelInfo represents a local BitNet model file
type ModelInfo struct {
    ID          string    `json:"id"`           // Unique identifier (usually filename)
    Name        string    `json:"name"`         // Display name
    Filename    string    `json:"filename"`     // Actual file name on disk
    FilePath    string    `json:"filepath"`     // Full absolute path
    Size        int64     `json:"size"`         // File size in bytes
    Modified    time.Time `json:"modified"`     // Last modified date
    IsDownloads bool      `json:"is_download"`  // True if currently downloading
}

// The DownloadStatus struct has been removed from this file to resolve the "redeclared" error.
// The single, authoritative definition for DownloadStatus is now in downloader.go.