package main

import (
	"log"

	"github.com/s19835/pg-opt-toolkit/internal/config"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "pgopt",
	Short: "PostgreSQL Query Optimizer Toolkit",
	Long:  "A CLI tool for analyzing and optimizing PostgreSQL queries",
}

var queryAnalzeCommand = &cobra.Command{
	Use:   "analyze [query]",
	Short: "Analyze a PostgreSQL query",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		analyzeQuery(query)
	},
}

func init() {
	rootCommand.AddCommand(queryAnalzeCommand)
}

func analyzeQuery(query string) {
	// load config
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Fail to load config: %v", err)
	}

	log.Println(query, cfg.URL)
}

func main() {
	analyzeQuery("query")
}
