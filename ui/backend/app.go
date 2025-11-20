package backend

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/mibrahimzia/bitnet-runner/internal/engine"
	"github.com/mibrahimzia/bitnet-runner/internal/embedder"
	"github.com/mibrahimzia/bitnet-runner/internal/models"
)

// App struct
type App struct {
	ctx          context.Context
	modelManager *models.Manager
	executor     *engine.Executor
	cancelFunc   context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		modelManager: models.NewManager(),
	}
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("App starting up...")

	
	// 1. Extract Engine on startup
	binPath, err := embedder.ExtractEngine()
	if err != nil {
		runtime.EventsEmit(ctx, "error", "Failed to load engine: "+err.Error())
		return
	}
	
	// 2. Initialize Executor with the path
	a.executor = engine.NewExecutor(binPath)
}

// shutdown is called at termination
func (a *App) Shutdown(ctx context.Context) {
	if a.executor != nil {
		a.executor.Shutdown()
	}
}

// -- Exposed Methods (Callable from JS) --

// ListModels returns the available models
func (a *App) ListModels() []models.ModelInfo {
	list, _ := a.modelManager.List()
	return list
}

// DownloadModel triggers a download and emits events for progress
func (a *App) DownloadModel(url string, name string) string {
	go func() {
		ch, err := a.modelManager.Download(url, name)
		if err != nil {
			runtime.EventsEmit(a.ctx, "download_error", err.Error())
			return
		}
		
		for status := range ch {
			// Emit event to Frontend: "download_progress"
			runtime.EventsEmit(a.ctx, "download_progress", status)
		}
	}()
	return "Download started"
} // Added missing closing brace

// LoadModelOnly starts the engine without chatting
func (a *App) LoadModelOnly(modelFile string) string {
	// 1. Resolve Path
	list, _ := a.modelManager.List()
	var fullPath string
	for _, m := range list {
		if m.ID == modelFile {
			fullPath = m.FilePath
			break
		}
	}
	if fullPath == "" {
		return "Error: Model not found"
	}

	// 2. Load
	go func() {
		err := a.executor.LoadModel(fullPath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "model_load_error", err.Error())
		} else {
			runtime.EventsEmit(a.ctx, "model_loaded", modelFile)
		}
	}()

	return "Loading started..."
}

// StartChat starts the inference and emits tokens via events
// StartChat starts the inference
func (a *App) StartChat(prompt string, modelFile string, temp float64, system string, topP float64, topK int, maxTokens int) string {
	// 1. Resolve Model Path
	list, _ := a.modelManager.List()
	var fullPath string
	for _, m := range list {
		if m.ID == modelFile {
			fullPath = m.FilePath
			break
		}
	}
	
	if fullPath == "" {
		return "Error: Model not found"
	}

	// 2. Config
	cfg := engine.InferenceConfig{
		ModelPath:    fullPath,
		Prompt:       prompt,
		SystemPrompt: system,     // Use user value
		Temperature:  temp,       // Use user value
		TopP:         topP,       // Use user value
		TopK:         topK,       // Use user value
		MaxTokens:    maxTokens,  // Use user value
		RepeatPenalty: 1.1,       // Fixed for now, or add argument
		Threads:      4,
	}

	// 3. Run in background
	go func() {
		stream, err := a.executor.StartInference(cfg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "chat_error", err.Error())
			return
		}

		for token := range stream {
			runtime.EventsEmit(a.ctx, "chat_token", token)
		}
		
		runtime.EventsEmit(a.ctx, "chat_done", true)
	}()

	return "Inference started"
}

// StopChat kills the current process
func (a *App) StopChat() {
	if a.executor != nil {
		a.executor.Stop()
	}
}