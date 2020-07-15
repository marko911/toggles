package evaluate

import (
	"reflect"
	"testing"
	"toggle/server/pkg/mock"
	"toggle/server/pkg/models"
)

type fields struct {
	r Repository
}
type args struct {
	e EvaluationData
}

type testCase struct {
	name    string
	fields  fields
	args    args
	want    *models.EvaluationResult
	wantErr bool
}

func Test_service_Evaluate(t *testing.T) {

	tests := []testCase{
		{
			name:   "flag with user targeting",
			fields: fields{r: &mock.Evaluate{FlagPath: "../../static/flagUserTargets.json"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				r: tt.fields.r,
			}
			got, err := s.Evaluate(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
