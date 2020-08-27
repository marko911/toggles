package evaluate

import (
	"errors"
	"fmt"
	er "toggle/server/pkg/errors"
	"toggle/server/pkg/models"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type userKey struct {
	Key string `json:"key"`
}

const percentMultiplier float64 = 100

// MatchFlagTarget parses all flag rules returning a variation if
// user matches
func (e *EvaluationRequest) MatchFlagTarget(flag *models.Flag) (*models.Variation, error) {

	for _, target := range flag.Targets {

		m, ok := e.User.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Cannot read evaluation request properly %v", m)
		}
		attrs := m["attributes"].(map[string]interface{})
		var u userKey
		err := mapstructure.Decode(e.User, &u)
		if err != nil {
			return nil, errors.New("Failed decoding user from evaluation request object")
		}

		var userMatches = true

		// SegmentMatch will return true if there is no segment
		if target.HasSegment() {
			fmt.Println("HAS SEGM")
			segMatch, err := target.SegmentMatch(attrs, u.Key)
			if err != nil {
				fmt.Println("Errrrrrrrrrr", err)
				return nil, err
			}
			userMatches = userMatches && segMatch
		}

		if target.HasRules() {
			fmt.Println("HAS RULES")
			// pass request data into expression evaluator
			rulesMatch, err := target.RulesMatch(attrs)
			if err != nil {
				logrus.Error("Error getting expression from target ", err)
			}
			fmt.Println("RULES MATCH", rulesMatch)
			userMatches = userMatches && rulesMatch
		}
		fmt.Println("USER MATCHES", userMatches)
		if userMatches {

			if target.HasRolloutDistribution() {

				fraction := CohortFraction(fmt.Sprintf("%s-%s", u.Key, flag.Key))
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
func (e *EvaluationRequest) MatchDefaultVariations(f models.Flag) (*models.Variation, error) {
	var u userKey
	err := mapstructure.Decode(e.User, &u)
	if err != nil {
		return nil, errors.New("Failed decoding user from evaluation request object")
	}

	fraction := CohortFraction(fmt.Sprintf("%s-%s", u.Key, f.Key))
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
