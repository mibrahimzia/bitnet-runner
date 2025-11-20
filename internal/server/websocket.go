package server

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

// Add these type definitions:
type ChatRequest struct {
    Message string `json:"message"`
    Model   string `json:"model"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

var upgrader = websocket.Upgrader{
    // Allow connections from any origin (required for local dev between ports)
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// WsHandler upgrades the HTTP connection to a WebSocket
func (s *Server) WsHandler(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // 1. Read the initial configuration from the socket
    // The UI sends the ChatRequest JSON immediately after connecting
    var req ChatRequest
    if err := conn.ReadJSON(&req); err != nil {
        conn.WriteJSON(ErrorResponse{Error: "Invalid request format"})
        return
    }

    // 2. Validate model availability
    // (In a real app, check if req.Model exists via s.modelManager)

    // 3. Configure the engine
    // We need to map API request to Engine config
    // This assumes you will add a helper to map types, but for now we do it manually
    // engineConfig := engine.InferenceConfig{...} (We will link this in the handler step)
    
    // For this step, we just set up the structure. 
    // The actual linking happens in the next file (handlers.go) to keep logic clean.
}