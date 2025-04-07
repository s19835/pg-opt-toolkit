package models

// Database configuration
type PGConfig struct {
	URL string // postgres://user:pass@host:port/db?sslmode=disable
}

// Single Node representation in query
type PlanNode struct {
	NodeType     string      `json:"Node Type"`
	RelationName string      `json:"Relation Name"`
	Alias        string      `json:"Alias"`
	StartupCost  float64     `json:"Startup Cost"`
	TotalCost    float64     `json:"Total Cost"`
	PlanRows     int64       `json:"Plan Rows"`
	PlanWidth    int64       `json:"Plan Width"`
	ActualTime   float64     `json:"Actual Total Time"`
	ActualRows   int64       `json:"Actual Rows"`
	Loops        int64       `json:"Actual Loops"`
	Filter       string      `json:"Filter,omitempty"`
	Plans        []*PlanNode `json:"Plans,omitempty"`
}

type OptimizationSavings struct {
	TimeSaved     float64 //milliseconds
	CostReduction float64 // unit cost
	RowsProcessed int64   // rows affected
}
