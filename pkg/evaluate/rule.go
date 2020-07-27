package evaluate

import (
	"errors"
	"fmt"
	er "toggle/server/pkg/errors"
	"toggle/server/pkg/models"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/zhouzhuojie/conditions"
)

type user struct {
	Key string `json:"key"`
}

type userAttr struct {
	Key        string
	Attributes map[string]interface{}
}

const percentMultiplier float64 = 100

// MatchFlagTarget parses all flag rules returning a variation if
// user matches
func (e *EvaluationData) MatchFlagTarget(targets []models.Target) (*models.Variation, error) {

	for _, target := range targets {

		m, ok := e.User.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Cannot read evaluation request properly %v", m)
		}
		attrs := m["attributes"].(map[string]interface{})

		// build expression of target
		expr, err := target.ToExpr()

		if err != nil {
			logrus.Error("Error getting expression from target ", err)
			return nil, err
		}
		// pass request data into expression evaluator
		match, err := conditions.Evaluate(expr, attrs)

		if err != nil {
			return nil, err
		}

		if match {

			if target.HasRolloutDistribution() {
				var u user
				err := mapstructure.Decode(e.User, &u)
				if err != nil {
					return nil, errors.New("Failed decoding user from evaluation request object")
				}

				fraction := CohortFraction(fmt.Sprintf("%s-%s", u.Key, e.FlagKey))
				percents := Percents(target.Variations)
				var inRange bool
				var min, accumulatedMax float64

				for i, p := range percents {
					if i == 0 {
						min = 0
					} else {
						min = accumulatedMax
					}
					accumulatedMax += p / percentMultiplier

					inRange = InRolloutRange(min, accumulatedMax, fraction)
					if inRange {
						variation := target.Variations[i]
						return &variation, nil
					}
				}

			}
			return target.GetMatchingVariation(), nil
		}
	}
	// no rule matched
	return nil, nil
}

// MatchDefaultVariations returns the default variation for this user
func (e *EvaluationData) MatchDefaultVariations(f *models.Flag) (*models.Variation, error) {
	var u user
	err := mapstructure.Decode(e.User, &u)
	if err != nil {
		return nil, errors.New("Failed decoding user from evaluation request object")
	}

	fraction := CohortFraction(fmt.Sprintf("%s-%s", u.Key, e.FlagKey))
	percents := Percents(f.Variations)
	var inRange bool
	var min, accumulatedMax float64

	for i, p := range percents {
		if i == 0 {
			min = 0
		} else {
			min = accumulatedMax
		}
		accumulatedMax += p / percentMultiplier

		inRange = InRolloutRange(min, accumulatedMax, fraction)
		if inRange {
			variation := f.Variations[i]
			return &variation, nil
		}
	}

	return nil, er.ErrVariationNotFound

}

// InRolloutRange checks if a value is in a range
func InRolloutRange(min, max, val float64) bool {
	if val >= min && val < max {
		return true
	}

	return false
}

// Percents return a list of percentages from variations
func Percents(v []models.Variation) []float64 {
	percentages := []float64{}
	for _, vari := range v {
		percentages = append(percentages, vari.Percent)
	}
	return percentages
}
