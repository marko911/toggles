package models

import "fmt"

const (

	// ConstraintOperatorEQ captures enum value "EQ"
	ConstraintOperatorEQ string = "EQ"

	// ConstraintOperatorNEQ captures enum value "NEQ"
	ConstraintOperatorNEQ string = "NEQ"

	// ConstraintOperatorLT captures enum value "LT"
	ConstraintOperatorLT string = "LT"

	// ConstraintOperatorLTE captures enum value "LTE"
	ConstraintOperatorLTE string = "LTE"

	// ConstraintOperatorGT captures enum value "GT"
	ConstraintOperatorGT string = "GT"

	// ConstraintOperatorGTE captures enum value "GTE"
	ConstraintOperatorGTE string = "GTE"

	// ConstraintOperatorEREG captures enum value "EREG"
	ConstraintOperatorEREG string = "EREG"

	// ConstraintOperatorNEREG captures enum value "NEREG"
	ConstraintOperatorNEREG string = "NEREG"

	// ConstraintOperatorIN captures enum value "IN"
	ConstraintOperatorIN string = "IN"

	// ConstraintOperatorNOTIN captures enum value "NOTIN"
	ConstraintOperatorNOTIN string = "NOTIN"

	// ConstraintOperatorCONTAINS captures enum value "CONTAINS"
	ConstraintOperatorCONTAINS string = "CONTAINS"

	// ConstraintOperatorNOTCONTAINS captures enum value "NOTCONTAINS"
	ConstraintOperatorNOTCONTAINS string = "NOTCONTAINS"
)

// OperatorToExprMap maps from model to expression operator
var OperatorToExprMap = map[string]string{
	ConstraintOperatorEQ:          "==",
	ConstraintOperatorNEQ:         "!=",
	ConstraintOperatorLT:          "<",
	ConstraintOperatorLTE:         "<=",
	ConstraintOperatorGT:          ">",
	ConstraintOperatorGTE:         ">=",
	ConstraintOperatorEREG:        "=~",
	ConstraintOperatorNEREG:       "!~",
	ConstraintOperatorIN:          "IN",
	ConstraintOperatorNOTIN:       "NOT IN",
	ConstraintOperatorCONTAINS:    "CONTAINS",
	ConstraintOperatorNOTCONTAINS: "NOT CONTAINS",
}

func (r *Rule) toExprStr() (string, error) {
	if r.Attribute == "" || r.Operator == "" || r.Value == "" {
		return "", fmt.Errorf(
			"empty Attribute/Operator/Value: %s/%s/%s",
			r.Attribute,
			r.Operator,
			r.Value,
		)
	}
	o, ok := OperatorToExprMap[r.Operator]
	if !ok {
		return "", fmt.Errorf("not supported operator: %s", r.Operator)
	}

	return fmt.Sprintf("({%s} %s %s)", r.Attribute, o, r.Value), nil
}
