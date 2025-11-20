package engine

// InferenceConfig holds parameters for text generation
type InferenceConfig struct {
	ModelPath     string  `json:"model_path"`
	Prompt        string  `json:"prompt"`
	SystemPrompt  string  `json:"system_prompt"`
	Temperature   float64 `json:"temperature"`     // 0.0 to 2.0
	TopP          float64 `json:"top_p"`           // 0.0 to 1.0
	TopK          int     `json:"top_k"`           // 0 to 100
	RepeatPenalty float64 `json:"repeat_penalty"`  // 1.0 to 2.0
	MaxTokens     int     `json:"max_tokens"`      // -1 for infinite
	Threads       int     `json:"threads"`         // number of CPU threads
}

func DefaultConfig() InferenceConfig {
	return InferenceConfig{
		Temperature:   0.8,
		TopP:          0.9,
		TopK:          40,
		RepeatPenalty: 1.1,
		MaxTokens:     512,
		Threads:       4,
	}
}