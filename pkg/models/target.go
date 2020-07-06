package models

import (
	"fmt"
	"strings"

	"github.com/zhouzhuojie/conditions"
)

// Target is a specific user constraint
type Target struct {
	Rules      []Rule      `json:"rule" bson:"rule"`             // slice used to allow for multiple rules to be used as an AND condition
	Variations []Variation `json:"variations" bson:"variations"` // distribution of variations if all rules pass
}

// ToExpr maps ConstraintArray to expr by joining 'AND'
func (t Target) ToExpr() (conditions.Expr, error) {
	strs := make([]string, 0, len(t.Rules))
	for _, c := range t.Rules {
		s, err := c.toExprStr()
		if err != nil {
			return nil, err
		}
		strs = append(strs, s)
	}
	exprStr := strings.Join(strs, " AND ")
	fmt.Println("EXPRESSSNs", exprStr)
	p := conditions.NewParser(strings.NewReader(exprStr))
	expr, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("%s. Note: if it's string or array of string, wrap it with quotes \"...\"", err)
	}
	return expr, nil
}

// HasRolloutDistribution checks target variations and returns weather
// it has a percent rollout for each or is a simple on/off
func (t Target) HasRolloutDistribution() bool {
	for _, v := range t.Variations {
		if v.Percent > 0 && v.Percent < 100 {
			return true
		}
	}
	return false
}

// GetMatchingVariation returns the ON variation of a simple bool flag
func (t Target) GetMatchingVariation() *Variation {
	for _, v := range t.Variations {
		if v.Percent == 100 {
			return &v
		}
	}
	return nil
}

// Percentages returns a slice of varations' percentages
func (t Target) Percentages() []float64 {
	percentages := []float64{}
	for _, v := range t.Variations {
		percentages = append(percentages, v.Percent)
	}
	return percentages
}

// // Rollout returns a variation based on rollout percentages
// func (t Target) Rollout(salt string) *Variation {
// 	percents := t.Percentages()
// 	// return something else

// 	return nil
// }
