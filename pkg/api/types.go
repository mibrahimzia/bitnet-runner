package api

// ChatRequest is the payload sent by the UI to start generation
type ChatRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	System      string  `json:"system_prompt"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	TopK        int     `json:"top_k"`
	MaxTokens   int     `json:"max_tokens"`
	Stream      bool    `json:"stream"` // If true, use WebSocket
}

// ChatResponse is a single chunk of generated text
type ChatResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// ModelDownloadRequest triggers a new download
type ModelDownloadRequest struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

// ErrorResponse is a standard error wrapper
type ErrorResponse struct {
	Error string `json:"error"`
}