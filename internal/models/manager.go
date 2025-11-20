package models

import (
	"sync"
)

type Manager struct {
	mu sync.Mutex
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) List() ([]ModelInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return ScanModels()
}

func (m *Manager) Download(url string, name string) (<-chan DownloadStatus, error) {
	// Create a channel for updates
	ch := make(chan DownloadStatus, 100)

	// Run download in background
	go func() {
		defer close(ch)
		_, err := DownloadModel(url, name, ch)
		if err != nil {
			ch <- DownloadStatus{
				ModelName: name,
				Error:     err.Error(),
			}
		}
	}()

	return ch, nil
}