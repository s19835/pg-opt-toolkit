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

	sb.WriteString(fmt.Sprintf("Execution Time: %.2f", plan.ExecutionTime))
	sb.WriteString(fmt.Sprintf("\nPlanning Time: %.2f", plan.Planning.PlanningTime))
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
	// Simple heuristic: nodes taking more than 100ms or with high relative cost
	if node.ActualTime > 100 {
		*bottlenecks = append(*bottlenecks, fmt.Sprintf("Slow operation: %s (%.2f ms)", node.NodeType, node.ActualTime))
	}

	for _, child := range node.Plans {
		a.findBottlenecks(child, bottlenecks)
	}
}
