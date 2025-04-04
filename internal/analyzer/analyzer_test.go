package analyzer_test

import (
	"encoding/json"
	"testing"

	"github.com/s19835/pg-opt-toolkit/internal/analyzer"
	"github.com/s19835/pg-opt-toolkit/pkg/models"
	"github.com/s19835/pg-opt-toolkit/pkg/queryplan"
	"github.com/stretchr/testify/assert"
)

const complexPlanJSON = `{
	"Plan": {
		"Node Type": "Limit",
		"Startup Cost": 1024.58,
		"Total Cost": 1026.12,
		"Plan Rows": 100,
		"Plan Width": 128,
		"Actual Startup Time": 1254.32,
		"Actual Total Time": 1256.78,
		"Actual Rows": 100,
		"Actual Loops": 1,
		"Plans": [
			{
				"Node Type": "Sort",
				"Startup Cost": 1024.58,
				"Total Cost": 1030.45,
				"Plan Rows": 2345,
				"Plan Width": 128,
				"Actual Startup Time": 1254.30,
				"Actual Total Time": 1255.12,
				"Actual Rows": 2345,
				"Actual Loops": 1,
				"Sort Key": ["users.created_at DESC"],
				"Plans": [
					{
						"Node Type": "Hash Join",
						"Startup Cost": 25.88,
						"Total Cost": 924.58,
						"Plan Rows": 2345,
						"Plan Width": 128,
						"Actual Startup Time": 45.32,
						"Actual Total Time": 123.45,
						"Actual Rows": 10000,
						"Actual Loops": 1,
						"Hash Cond": "(users.id = orders.user_id)",
						"Plans": [
							{
								"Node Type": "Seq Scan",
								"Parent Relationship": "Outer",
								"Relation Name": "users",
								"Alias": "users",
								"Startup Cost": 0.00,
								"Total Cost": 125.88,
								"Plan Rows": 10000,
								"Plan Width": 64,
								"Actual Startup Time": 5.32,
								"Actual Total Time": 25.45,
								"Actual Rows": 10000,
								"Actual Loops": 1,
								"Filter": "(status = 'active')"
							},
							{
								"Node Type": "Hash",
								"Parent Relationship": "Inner",
								"Startup Cost": 15.88,
								"Total Cost": 15.88,
								"Plan Rows": 500,
								"Plan Width": 64,
								"Actual Startup Time": 15.32,
								"Actual Total Time": 15.32,
								"Actual Rows": 500,
								"Actual Loops": 1,
								"Plans": [
									{
										"Node Type": "Seq Scan",
										"Relation Name": "orders",
										"Alias": "orders",
										"Startup Cost": 0.00,
										"Total Cost": 15.88,
										"Plan Rows": 500,
										"Plan Width": 64,
										"Actual Startup Time": 0.32,
										"Actual Total Time": 5.45,
										"Actual Rows": 500,
										"Actual Loops": 1,
										"Filter": "(created_at > (now() - '30 days'::interval))"
									}
								]
							}
						]
					}
				]
			}
		]
	},
	"Planning Time": 12.34,
	"Execution Time": 1268.12
}`

func TestPerformanceMatrixCalculation(t *testing.T) {
	var plan queryplan.QueryPlan
	err := json.Unmarshal([]byte(complexPlanJSON), &plan)
	assert.NoError(t, err)

	analyzer := analyzer.NewQueryAnalyzer()
	result, err := analyzer.Analyze(&plan)
	assert.NoError(t, err)

	// Test latency measurement
	assert.Contains(t, result, "Execution Time: 1268.12 ms")
	assert.Contains(t, result, "Planning Time: 12.34 ms")

	// Test const measurement
	assert.Contains(t, result, "Cost: 1024.58..1026.12")
	assert.Contains(t, result, "Cost: 25.88..924.58")
}

func TestBottlenecksDetection(t *testing.T) {
	var plan queryplan.QueryPlan
	err := json.Unmarshal([]byte(complexPlanJSON), &plan)
	assert.NoError(t, err)

	analyzer := analyzer.NewQueryAnalyzer()
	bottlenecks := analyzer.IdentifyBottlenecks(&plan)

	expectedBottlenecks := []string{
		"Slow operation: Sort (1255.12 ms)",
		"Slow operation: Hash Join (123.45 ms)",
		"Inefficient operation: Seq Scan on users (25.45 ms for 10000 rows)",
		"Potential optimization: Filter on status could use index",
		"Potential optimization: Filter on created_at could use index",
	}

	assert.Len(t, bottlenecks, len(expectedBottlenecks))
	for _, expected := range expectedBottlenecks {
		assert.Contains(t, bottlenecks, expected)
	}
}

func TestCostSavingsEstimation(t *testing.T) {
	var plan queryplan.QueryPlan
	err := json.Unmarshal([]byte(complexPlanJSON), &plan)
	assert.NoError(t, err)

	analyzer := analyzer.NewQueryAnalyzer()

	// test seq scan node
	seqScanNode := findNode(&plan.Plan, "Seq Scan", "users")
	assert.NotNil(t, seqScanNode)

	estimatedIndexScan := &models.PlanNode{
		NodeType:    "Index Scan",
		StartupCost: 5.00,
		TotalCost:   50.00,
		ActualTime:  2.50,
	}

	savings := analyzer.EstimateSavings(seqScanNode, estimatedIndexScan)
	assert.Greater(t, savings.TimeSaved, 20.0)     //should be > 20ms
	assert.Greater(t, savings.CostReduction, 50.0) //should be > 50 cost unit
}

func findNode(node *models.PlanNode, nodeType string, relation string) *models.PlanNode {
	if node.NodeType == nodeType && (relation == "" || node.RelationName == relation) {
		return node
	}

	for _, child := range node.Plans {
		if found := findNode(child, nodeType, relation); found != nil {
			return found
		}
	}
	return nil
}
