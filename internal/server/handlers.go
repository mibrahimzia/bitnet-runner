package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mibrahimzia/bitnet-runner/internal/engine"
	//"github.com/mibrahimzia/bitnet-runner/internal/models"
	"github.com/mibrahimzia/bitnet-runner/pkg/api"
)

// HandleListModels returns all available local models
func (s *Server) HandleListModels(c *gin.Context) {
	models, err := s.modelManager.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models)
}

// HandlePullModel triggers a background download
func (s *Server) HandlePullModel(c *gin.Context) {
	var req api.ModelDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Start download (non-blocking)
	// Note: In a full implementation, you'd store the channel in a map to track progress via another endpoint
	_, err := s.modelManager.Download(req.Url, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "download_started", "model": req.Name})
}

// HandleChatStream manages the WebSocket connection and Engine execution
func (s *Server) HandleChatStream(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var req api.ChatRequest
	if err := conn.ReadJSON(&req); err != nil {
		return 
	}

	// Prepare Engine Config
	cfg := engine.InferenceConfig{
		ModelPath:     req.Model, // This should be the full path resolved by manager
		Prompt:        req.Prompt,
		SystemPrompt:  req.System,
		Temperature:   req.Temperature,
		TopP:          req.TopP,
		TopK:          req.TopK,
		MaxTokens:     req.MaxTokens,
		Threads:       4, // Default
	}

	// Initialize Executor
	// Note: We resolve the binary path dynamically
	binPath, _ := s.getBinaryPath() 
	exec := engine.NewExecutor(binPath)

	// Start Inference
	stream, err := exec.StartInference(cfg)
	if err != nil {
		conn.WriteJSON(api.ErrorResponse{Error: err.Error()})
		return
	}

	// Stream loop
	for token := range stream {
		resp := api.ChatResponse{
			Content: token,
			Done:    false,
		}
		if err := conn.WriteJSON(resp); err != nil {
			exec.Stop()
			break
		}
	}

	// Send done signal
	conn.WriteJSON(api.ChatResponse{Done: true})
}