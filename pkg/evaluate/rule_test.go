package evaluate

import (
	"testing"
	"toggle/server/pkg/models"

	"github.com/dchest/uniuri"
	"github.com/google/go-cmp/cmp"
)

var randomKeyCharset = []byte("123456789abcdefghijkmnopqrstuvwxyz")

func GenerateFlag(a, o, v string) *models.Flag {
	targets := []models.Target{
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

	return &models.Flag{
		Targets: targets,
	}
}

func TestEvaluationRequest_MatchFlagTarget(t *testing.T) {
	type fields struct {
		User interface{}
	}
	type args struct {
		flag *models.Flag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Variation
		wantErr bool
	}{
		{
			"CONTAINS operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"groups": []string{"males"},
				},
			},
			},
			args{flag: GenerateFlag("groups", models.ConstraintOperatorCONTAINS, "\"males\"")}, // always have to provide quoted value to evaluation lib
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"NOT CONTAINS operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"groups": []string{"gorillas"},
				},
			},
			},
			args{flag: GenerateFlag("groups", models.ConstraintOperatorNOTCONTAINS, "\"males\"")}, // always have to provide quoted value to evaluation lib
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"IN operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"name": "marko",
				},
			},
			},
			args{flag: GenerateFlag("name", models.ConstraintOperatorIN, `["marko","hana"]`)},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"NOT IN operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"name": "jack",
				},
			},
			},
			args{flag: GenerateFlag("name", models.ConstraintOperatorNOTIN, `["marko","hana"]`)},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"EQ operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"name": "jack",
				},
			},
			},
			args{flag: GenerateFlag("name", models.ConstraintOperatorEQ, "\"jack\"")},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"LT operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"age": 25,
				},
			},
			},
			args{flag: GenerateFlag("age", models.ConstraintOperatorLT, "40")},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"LTE operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"age": 40,
				},
			},
			},
			args{flag: GenerateFlag("age", models.ConstraintOperatorLTE, "40")},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
		{
			"EREG operator", fields{map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"regex": "John",
				},
			},
			},
			args{flag: GenerateFlag("regex", models.ConstraintOperatorEREG, `"([A-Z])\w+"`)},
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
	}

	t.Log("Given the need to use operators on user request attributes")

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			e := &EvaluationRequest{
				User: tt.fields.User,
			}
			got, err := e.MatchFlagTarget(tt.args.flag)

			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluationRequest.MatchFlagTarget() error = %v, wantErr %v", err, tt.wantErr)
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
