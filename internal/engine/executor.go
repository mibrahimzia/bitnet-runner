package engine

import (
	"bufio"
	"bytes"
	"context" // Added context
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"syscall" 
	"time"
)

type ServerRequest struct {
	Prompt        string  `json:"prompt"`
	NPredict      int     `json:"n_predict"`
	Temperature   float64 `json:"temperature"`
	TopP          float64 `json:"top_p"`
	TopK          int     `json:"top_k"`
	RepeatPenalty float64 `json:"repeat_penalty"`
	Stream        bool    `json:"stream"`
}

type ServerResponse struct {
	Content string `json:"content"`
	Stop    bool   `json:"stop"`
}

type Executor struct {
	binPath     string
	cmd         *exec.Cmd
	mu          sync.Mutex
	running     bool
	activeModel string
	serverPort  string
	
	// Context for the active chat request
	cancelRequest context.CancelFunc
}

func NewExecutor(binaryPath string) *Executor {
	return &Executor{
		binPath:    binaryPath,
		serverPort: "8080",
	}
}

// LoadModel starts the server without running inference
func (e *Executor) LoadModel(modelPath string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running && e.activeModel == modelPath {
		return nil // Already loaded
	}

	return e.restartServerLocked(modelPath)
}

func (e *Executor) StartInference(config InferenceConfig) (<-chan string, error) {
	e.mu.Lock()
	// Auto-load if not ready
	if !e.running || e.activeModel != config.ModelPath {
		if err := e.restartServerLocked(config.ModelPath); err != nil {
			e.mu.Unlock()
			return nil, err
		}
	}
	
	// Cancel any previous request just in case
	if e.cancelRequest != nil {
		e.cancelRequest()
	}
	
	// Create new cancellable context for this specific request
	ctx, cancel := context.WithCancel(context.Background())
	e.cancelRequest = cancel
	e.mu.Unlock()

	// Use the simpler template from your working config
	fullPrompt := fmt.Sprintf("System: %sUser: %s<|eot_id|>Assistant:", config.SystemPrompt, config.Prompt)



	
	reqBody := ServerRequest{
		Prompt:        fullPrompt,
		NPredict:      config.MaxTokens,
		Temperature:   config.Temperature,
		TopP:          config.TopP,
		TopK:          config.TopK,
		RepeatPenalty: config.RepeatPenalty,
		Stream:        true,
	}

	jsonData, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("http://127.0.0.1:%s/completion", e.serverPort)
	
	// Create request with Context
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	outputChan := make(chan string)

	go func() {
		defer close(outputChan)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				jsonStr := strings.TrimPrefix(line, "data: ")
				if strings.Contains(jsonStr, "[DONE]") {
					return
				}
				var data ServerResponse
				if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
					outputChan <- data.Content
					if data.Stop {
						return
					}
				}
			}
		}
	}()

	return outputChan, nil
}

// Stop cancels the CURRENT GENERATION but keeps the server running
func (e *Executor) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if e.cancelRequest != nil {
		e.cancelRequest() // This cuts the HTTP connection
		e.cancelRequest = nil
	}
	return nil
}

// Shutdown actually kills the server process (used when closing app)
func (e *Executor) Shutdown() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if e.cancelRequest != nil {
		e.cancelRequest()
	}
	
	if e.cmd != nil && e.cmd.Process != nil {
		e.running = false
		return e.cmd.Process.Kill()
	}
	return nil
}

func (e *Executor) restartServerLocked(modelPath string) error {
	if e.cmd != nil && e.cmd.Process != nil {
		_ = e.cmd.Process.Kill()
		e.cmd.Wait()
	}

	args := []string{
		"-m", modelPath,
		"--port", e.serverPort,
		"-c", "2048",
		"--host", "127.0.0.1",
	}

	// fmt.Printf("DEBUG: Starting Server: %s %v\n", e.binPath, args) // Comment out debug log for production
	cmd := exec.Command(e.binPath, args...)
	
	// --- THIS HIDES THE BLACK WINDOW ---
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	// -----------------------------------
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
    // ... rest of function stays the same

	e.cmd = cmd
	e.running = true
	e.activeModel = modelPath

	// Wait for server health
	fmt.Println("DEBUG: Waiting for model to load...")
	for i := 0; i < 30; i++ { 
		time.Sleep(500 * time.Millisecond)
		_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/health", e.serverPort))
		if err == nil {
			fmt.Println("DEBUG: Model Loaded!")
			return nil
		}
	}
	return nil
}