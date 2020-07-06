package evaluate

import (
	"errors"
	"fmt"
	"toggle/server/pkg/models"

	"github.com/zhouzhuojie/conditions"
)

type User struct {
	Key string
}

const PercentMultiplier float64 = 10

// MatchFlagTarget parses all flag rules returning a variation if
// user matches
func (e *EvaluationData) MatchFlagTarget(targets []models.Target) (*models.Variation, error) {
	for _, target := range targets {

		m, ok := e.User.(map[string]interface{})
		if !ok {
			return nil, errors.New("Cannot read evaluation request properly")
		}

		// build expression of target
		expr, err := target.ToExpr()
		if err != nil {
			return nil, err
		}

		// pass request data into expression evaluator
		match, err := conditions.Evaluate(expr, m)
		if err != nil {
			return nil, err
		}

		if match {
			if target.HasRolloutDistribution() {
				u, ok := e.User.(User)
				if !ok {
					return nil, errors.New("Cannot read evaluation request properly")
				}
				fraction := CohortFraction(fmt.Sprintf("%s-%s", u.Key, e.FlagKey))
				//do hashing here and return variation
				percents := target.Percentages()

				var inRange bool

				// abstract this away?

				for i, p := range percents {
					if i == 0 {
						inRange = InRolloutRange(0, p/PercentMultiplier, fraction)
					} else {
						inRange = InRolloutRange(percents[i-1], p/PercentMultiplier, fraction)
					}
					if inRange {
						variation := target.Variations[i]
						return &variation, nil
					}
				}

			}
			return target.GetMatchingVariation(), nil
		}
	}
	// change
	return nil, nil
}

// InRolloutRange checks if a value is in a range
func InRolloutRange(min, max, val float64) bool {
	if val >= min && val < max {
		return true
	}

	return false
}
