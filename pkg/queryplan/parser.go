package queryplan

import (
	"encoding/json"
	"fmt"

	"github.com/s19835/pg-opt-toolkit/pkg/models"
)

type QueryPlan struct {
	Plan     models.PlanNode `json:"Plan"`
	Planning struct {
		PlanningTime float64 `json:"Planning Time"`
	} `json:"Planning"`
	ExecutionTime float64 `json:"Execution Time"`
}

func ParsePlanJSON(jsonData string) (*QueryPlan, error) {
	var plan []QueryPlan
	err := json.Unmarshal([]byte(jsonData), &plan)

	if err != nil {
		return nil, fmt.Errorf("failed to parse query plan JSON: %w", err)
	}

	if len(plan) == 0 {
		return nil, fmt.Errorf("empty query plan")
	}

	return &plan[0], nil
}
