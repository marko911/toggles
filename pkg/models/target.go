package models

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zhouzhuojie/conditions"
)

// Target is a specific user constraint
type Target struct {
	Rules []Rule `json:"rules" bson:"rules"` // slice used to allow for multiple rules to be used as an AND condition
	// Users      []string    `json:"users,omitempty" bson:"users,omitempty"` // user keys
	Variations []Variation `json:"variations" bson:"variations"` // distribution of variations if all rules pass
	Segment    `json:"segment,omitempty" bson:"segment,omitempty"`
}

// ToExpr maps ConstraintArray to expr by joining 'AND'
func (t Target) ToExpr() (conditions.Expr, error) {
	strs := make([]string, 0, len(t.Rules))
	for _, c := range t.Rules {
		//TODO: fix this hacky way of segmeents in rules
		if c.Attribute == "segment" {
			continue
		}
		s, err := c.toExprStr()
		if err != nil {
			return nil, err
		}
		strs = append(strs, s)
	}
	exprStr := strings.Join(strs, " AND ")
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

// HasSegment checks if target is defined with a segment
func (t Target) HasSegment() bool {
	return t.Segment.ID != ""
}

// HasRules checks if target is defined with rules
func (t Target) HasRules() bool {
	var rules []Rule
	for _, r := range t.Rules {
		if r.Attribute != "segment" {
			rules = append(rules, r)
		}
	}
	return len(rules) > 0
}

// SegmentMatch tells whether user fits segment or not
func (t Target) SegmentMatch(attrs map[string]interface{}, userKey string) (bool, error) {
	if t.Segment.ID == "" {
		return true, nil
	}
	//first check if user is in segment users list
	_, found := Find(t.Segment.Users, userKey)
	if found {
		return true, nil
	}

	//now parse through segment rules
	expr, err := t.Segment.ToExpr()
	if err != nil {
		logrus.Error("Error getting expression from target ", err)
		return false, err
	}

	return conditions.Evaluate(expr, attrs)

}

//RulesMatch tells us whether rules matched user
func (t Target) RulesMatch(attrs map[string]interface{}) (bool, error) {
	expr, err := t.ToExpr()
	if err != nil {
		logrus.Error("Error getting expression from target ", err)
		return false, err
	}
	// pass request data into expression evaluator
	return conditions.Evaluate(expr, attrs)

}

// Find looks for a string in a slice and returns its index and weather it was successful
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
