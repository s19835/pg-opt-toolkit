package analyzer

import (
	"fmt"
	"strings"

	"github.com/s19835/pg-opt-toolkit/pkg/models"
	"github.com/s19835/pg-opt-toolkit/pkg/queryplan"
)

type QueryAnalyzer struct{}

func NewQueryAnalyzer() *QueryAnalyzer {
	return &QueryAnalyzer{}
}

func (a *QueryAnalyzer) Analyze(plan *queryplan.QueryPlan) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Execution Time: %.2f ms", plan.ExecutionTime))
	sb.WriteString(fmt.Sprintf("\nPlanning Time: %.2f ms", plan.PlanningTime))
	sb.WriteString("\n")

	a.analyzeNode(&plan.Plan, &sb, 0)

	return sb.String(), nil
}

func (a *QueryAnalyzer) analyzeNode(node *models.PlanNode, sb *strings.Builder, depth int) {
	indent := strings.Repeat(" ", depth)

	sb.WriteString(fmt.Sprintf("%sNode Type: %s\n", indent, node.NodeType))

	if node.RelationName != "" {
		sb.WriteString(fmt.Sprintf("%sRelation Name: %s (alias: %s)\n", indent, node.RelationName, node.Alias))
	}
	sb.WriteString(fmt.Sprintf("%sCost: %.2f..%2f\n", indent, node.StartupCost, node.TotalCost))
	sb.WriteString(fmt.Sprintf("%sActual Time: %.2f ms, Rows: %d, Loops: %d\n", indent, node.ActualTime, node.ActualRows, node.Loops))
	sb.WriteString("\n")

	for _, child := range node.Plans {
		a.analyzeNode(child, sb, depth+1)
	}
}

func (a *QueryAnalyzer) IdentifyBottlenecks(plan *queryplan.QueryPlan) []string {
	var bottlenecks []string
	a.findBottlenecks(&plan.Plan, &bottlenecks)
	return bottlenecks
}

func (a *QueryAnalyzer) findBottlenecks(node *models.PlanNode, bottlenecks *[]string) {
	// Detect slow operations (over 100ms)
	if node.ActualTime > 100 {
		*bottlenecks = append(*bottlenecks,
			fmt.Sprintf("Slow operation: %s (%.2f ms)", node.NodeType, node.ActualTime))
	}

	// Detect inefficient scans
	if node.NodeType == "Seq Scan" {
		// Calculate rows processed per millisecond
		rowsPerMs := float64(node.ActualRows) / node.ActualTime

		relationInfo := ""
		if node.RelationName != "" {
			relationInfo = fmt.Sprintf(" on %s", node.RelationName)
		}

		// Flag as inefficient if < 50 rows/ms or if scanning large table (>1000 rows)
		if rowsPerMs < 50 || node.ActualRows > 1000 {
			*bottlenecks = append(*bottlenecks,
				fmt.Sprintf("Inefficient operation: Seq Scan%s (%.2f ms for %d rows)",
					relationInfo, node.ActualTime, node.ActualRows))
		}

		// Check for filter conditions that could use indexes
		if node.Filter != "" {
			*bottlenecks = append(*bottlenecks,
				fmt.Sprintf("Potential optimization: Filter on %s could use index", node.Filter))
		}
	}

	// Recursively check child nodes
	for _, child := range node.Plans {
		a.findBottlenecks(child, bottlenecks)
	}
}

// for testing and optimizing
func (a *QueryAnalyzer) EstimateSavings(current *models.PlanNode, optimized *models.PlanNode) *models.OptimizationSavings {
	return &models.OptimizationSavings{
		TimeSaved:     current.ActualTime - optimized.ActualTime,
		CostReduction: current.TotalCost - optimized.TotalCost,
		RowsProcessed: current.ActualRows,
	}
}

func (a *QueryAnalyzer) CalculateProcessingRate(rows int64, time float64, loop int64) float64 {
	if time <= 0 {
		return 0
	}
	return float64(rows*loop) / time
}
