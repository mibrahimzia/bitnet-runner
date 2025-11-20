package server

import (
	"fmt"
	//"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mibrahimzia/bitnet-runner/internal/embedder"
	"github.com/mibrahimzia/bitnet-runner/internal/models"
)

type Server struct {
	router       *gin.Engine
	modelManager *models.Manager
	port         string
	binPath      string // Path to the extracted bitnet.exe
}

func NewServer(port string) (*Server, error) {
	// 1. Initialize Model Manager
	mm := models.NewManager()

	// 2. Extract Engine Resources (ensure bitnet.exe is ready)
	binPath, err := embedder.ExtractEngine()
	if err != nil {
		return nil, fmt.Errorf("engine init failed: %w", err)
	}

	s := &Server{
		router:       gin.Default(),
		modelManager: mm,
		port:         port,
		binPath:      binPath,
	}

	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {
	// CORS configuration to allow UI to talk to localhost server
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	s.router.Use(cors.New(config))

	api := s.router.Group("/api/v1")
	{
		api.GET("/models", s.HandleListModels)
		api.POST("/models/pull", s.HandlePullModel)
		// WebSocket endpoint
		api.GET("/chat", s.HandleChatStream)
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.port)
}

// Helper to get binary path
func (s *Server) getBinaryPath() (string, error) {
	if s.binPath == "" {
		return "", fmt.Errorf("engine not initialized")
	}
	return s.binPath, nil
}