package models

import (
	"fmt"
	"strings"

	"github.com/zhouzhuojie/conditions"
)

// ToExpr maps ConstraintArray to expr by joining 'AND'
func (s Segment) ToExpr() (conditions.Expr, error) {
	strs := make([]string, 0, len(s.Rules))
	for _, c := range s.Rules {
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
