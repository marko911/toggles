package evaluate

import (
	"testing"
	"toggle/server/pkg/models"

	"github.com/dchest/uniuri"
	"github.com/google/go-cmp/cmp"
)

var randomKeyCharset = []byte("123456789abcdefghijkmnopqrstuvwxyz")

func NewTargetsWithGroup(a, v, o string) []models.Target {
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
			args{targets: NewTargetsWithGroup("groups", "\"males\"", "CONTAINS")}, // always have to provide quoted value to evaluation lib
			&models.Variation{
				Name: "On", Percent: 100,
			},
			false,
		},
	}

	t.Log("Given the need to use operators on user request attributes")

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Logf("When using the %s and providing a matching attribute", tt.name)

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
			t.Logf("\t\tShould receive a \"%s\" variation. %v",
				tt.want.Name, checkMark)
		})
	}
}
