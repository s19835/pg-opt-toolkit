package main

import (
	"fmt"
	"log"

	"github.com/s19835/pg-opt-toolkit/internal/analyzer"
	"github.com/s19835/pg-opt-toolkit/internal/config"
	"github.com/s19835/pg-opt-toolkit/internal/connector"
	"github.com/s19835/pg-opt-toolkit/pkg/models"
	"github.com/s19835/pg-opt-toolkit/pkg/queryplan"
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
		log.Fatalf("Failed to load config: %v", err)
	}

	// connect to database
	conn, err := connector.NewPGConnector(models.PGConfig{
		URL: cfg.URL,
	})
	if err != nil {
		log.Fatalf("Error connecting database: %v", err)
	}
	defer conn.Close()

	// analyze the provided query
	planJSON, err := conn.ExplainAnalyze(query)
	if err != nil {
		log.Fatalf("Unable to analyze query: %v", err)
	}

	// parse query plan
	plan, err := queryplan.ParsePlanJSON(planJSON)
	if err != nil {
		log.Fatalf("Failed to parse query: %v", err)
	}

	// analyze query
	analyzer := analyzer.NewQueryAnalyzer()
	analysis, err := analyzer.Analyze(plan)
	if err != nil {
		log.Fatalf("Failed to analyze query: %v", err)
	}

	fmt.Println("Query Analysis:")
	fmt.Println(analysis)

	// identify bottlenecks
	bottlenecks := analyzer.IdentifyBottlenecks(plan)
	if len(bottlenecks) > 0 {
		fmt.Println("\nPotential bottlenecks:")
		for _, b := range bottlenecks {
			fmt.Println("-", b)
		}
	}
}

func main() {
	analyzeQuery("SELECT * FROM urls")
}
