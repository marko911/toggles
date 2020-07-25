package evaluate

import (
	"testing"
	"toggle/server/pkg/models"

	"github.com/dchest/uniuri"
	"github.com/google/go-cmp/cmp"
)

var randomKeyCharset = []byte("123456789abcdefghijkmnopqrstuvwxyz")

func GenerateTarget(a, o, v string) []models.Target {
	r := []models.Target{
		{
			Rules: []models.Rule{
				{Attribute: a, Operator: o, Value: v},
			},
			Variations: []models.Variation{
				{Name: "On", Percent: 100},
				{Name: "Off", Percent: 0},
			},
		},
	}

	return r
}

func TestEvaluationData_MatchFlagTarget(t *testing.T) {
	type fields struct {
		FlagKey string
		User    interface{}
	}
	type args struct {
		targets []models.Target
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Variation
		wantErr bool
	}{
		{
			"CONTAINS operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"groups": []string{"males"},
				},
			},
			},
			args{targets: GenerateTarget("groups", models.ConstraintOperatorCONTAINS, "\"males\"")}, // always have to provide quoted value to evaluation lib
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"NOT CONTAINS operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"groups": []string{"gorillas"},
				},
			},
			},
			args{targets: GenerateTarget("groups", models.ConstraintOperatorNOTCONTAINS, "\"males\"")}, // always have to provide quoted value to evaluation lib
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"IN operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"name": "marko",
				},
			},
			},
			args{targets: GenerateTarget("name", models.ConstraintOperatorIN, `["marko","hana"]`)},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"NOT IN operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"name": "jack",
				},
			},
			},
			args{targets: GenerateTarget("name", models.ConstraintOperatorNOTIN, `["marko","hana"]`)},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"EQ operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"name": "jack",
				},
			},
			},
			args{targets: GenerateTarget("name", models.ConstraintOperatorEQ, "\"jack\"")},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"LT operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"age": 25,
				},
			},
			},
			args{targets: GenerateTarget("age", models.ConstraintOperatorLT, "40")},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"LTE operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"age": 40,
				},
			},
			},
			args{targets: GenerateTarget("age", models.ConstraintOperatorLTE, "40")},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"EREG operator", fields{"beta-users", map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"regex": "John",
				},
			},
			},
			args{targets: GenerateTarget("regex", models.ConstraintOperatorEREG, `"([A-Z])\w+"`)},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
	}

	t.Log("Given the need to use operators on user request attributes")

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			e := &EvaluationData{
				FlagKey: tt.fields.FlagKey,
				User:    tt.fields.User,
			}
			got, err := e.MatchFlagTarget(tt.args.targets)

			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluationData.MatchFlagTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Fatalf(diff)
			}
			t.Logf("\n\t When using the  %s , Should receive a \"%s\" variation. %v",
				tt.name, tt.want.Name, checkMark)
		})
	}
}
