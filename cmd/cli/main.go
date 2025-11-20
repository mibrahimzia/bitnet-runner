package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/mibrahimzia/bitnet-runner/internal/embedder"
	"github.com/mibrahimzia/bitnet-runner/internal/engine"
	"github.com/mibrahimzia/bitnet-runner/internal/models"
	"github.com/mibrahimzia/bitnet-runner/internal/server"
)

// Global flags
var (
	hostFlag string
	portFlag string
)

// Run flags
var (
	tempFlag   float64
	promptFlag string
)

var rootCmd = &cobra.Command{
	Use:   "bitnet",
	Short: "BitNet Inference Runner",
	Long:  `A local runner for BitNet b1.58 large language models.`,
}

func main() {
	// Define global flags
	rootCmd.PersistentFlags().StringVar(&hostFlag, "host", "localhost", "Server host")
	rootCmd.PersistentFlags().StringVar(&portFlag, "port", "8080", "Server port")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Register subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(runCmd)

	// Add flags to run command
	runCmd.Flags().Float64Var(&tempFlag, "temp", 0.8, "Temperature")
	runCmd.Flags().StringVarP(&promptFlag, "prompt", "p", "", "Prompt text")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Initializing BitNet Engine...\n")

		srv, err := server.NewServer(portFlag)
		if err != nil {
			fmt.Printf("Error initializing server: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Server listening on http://%s:%s\n", hostFlag, portFlag)
		if err := srv.Start(); err != nil {
			fmt.Printf("Server crashed: %v", err)
			os.Exit(1)
		}
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed models",
	Run: func(cmd *cobra.Command, args []string) {
		// Direct access to internal manager for CLI usage
		mgr := models.NewManager()
		list, err := mgr.List()
		if err != nil {
			fmt.Printf("Error listing models: %v", err)
			os.Exit(1)
		}

		if len(list) == 0 {
			fmt.Println("No models found in ~/.bitnet_runner/models/")
			return
		}

		fmt.Printf("%-25s %-15s %s\n", "MODEL ID", "SIZE", "PATH")
		fmt.Println("------------------------------------------------------------")
		for _, m := range list {
			sizeMB := m.Size / 1024 / 1024
			fmt.Printf("%-25s %d MB        %s\n", m.ID, sizeMB, m.FilePath)
		}
	},
}

var runCmd = &cobra.Command{
	Use:   "run [model_filename]",
	Short: "Run a model inference",
	Args:  cobra.ExactArgs(1), // Requires exactly 1 argument (the model name)
	Run: func(cmd *cobra.Command, args []string) {
		modelFile := args[0]

		// 1. Initialize Engine Resources
		binPath, err := embedder.ExtractEngine()
		if err != nil {
			fmt.Printf("Failed to extract engine: %v\n", err)
			os.Exit(1)
		}

		// 2. Resolve Model Path
		mgr := models.NewManager()
		list, _ := mgr.List()
		var fullPath string

		// Simple search
		for _, m := range list {
			if m.ID == modelFile || m.Filename == modelFile {
				fullPath = m.FilePath
				break
			}
		}

		if fullPath == "" {
			// Fallback: check if user provided a full path
			if _, err := os.Stat(modelFile); err == nil {
				fullPath = modelFile
			} else {
				fmt.Printf("Model '%s' not found.\n", modelFile)
				os.Exit(1)
			}
		}

		// 3. Prepare Config
		cfg := engine.InferenceConfig{
			ModelPath:   fullPath,
			Prompt:      promptFlag,
			Temperature: tempFlag,
			TopP:        0.9,
			TopK:        40,
			MaxTokens:   512,
			Threads:     4,
		}

		// If no prompt flag, read from stdin (simple interactive mode)
		if promptFlag == "" {
			fmt.Print("Enter prompt: ")
			fmt.Scanln(&cfg.Prompt)
		}

		fmt.Printf("Loading %s...\n", modelFile)

		// 4. Execute
		exec := engine.NewExecutor(binPath)
		stream, err := exec.StartInference(cfg)
		if err != nil {
			fmt.Printf("Error starting inference: %v\n", err)
			os.Exit(1)
		}

		fmt.Print("\nBitNet: ")
		for token := range stream {
			fmt.Print(token)
		}
		fmt.Println("\n")
	},
}