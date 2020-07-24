package evaluate

import (
	"testing"
	"toggle/server/pkg/mock"
	"toggle/server/pkg/models"

	"github.com/cheekybits/is"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/mgo.v2/bson"
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
	want    *EvaluationResult
	wantErr bool
}

const checkMark = "\u2713"
const ballotX = "\u2717"

func Test_service_Evaluate(t *testing.T) {

	tests := []testCase{
		{
			"flag with user targeting",
			fields{r: &mock.Evaluate{FlagPath: "../../static/flagUserTargets.json"}},
			args{EvaluationData{"hey-ladies", user{Key: "jenny@hey.com"}}},
			&EvaluationResult{
				models.User{Key: "jenny@hey.com"},
				&models.Variation{Name: "On", Percent: 100, UserKeys: []string{"jenny@hey.com", "mary@hey.com"}},
				bson.ObjectIdHex("5f09d08d40a5b800068a5d88"),
			},
			false,
		},
		{
			"default flag distribution",
		},
		// {
		// 	"correct distribution for targets on flag"
		// },
		// {
		// 	""
		// }
	}

	t.Log("Given the need to match a variation to a request:")
	{

		for _, tt := range tests {
			t.Logf("\tWhen checking \"%s\" for matching variation", tt.name)

			t.Run(tt.name, func(t *testing.T) {
				is := is.New(t)

				s := &service{
					r: tt.fields.r,
				}
				got, err := s.Evaluate(tt.args.e)
				is.NoErr(err)

				diff := cmp.Diff(tt.want, got)
				if diff != "" {
					t.Fatalf(diff)
				}
				t.Logf("\t\tShould receive a \"%s\" message. %v",
					tt.want.Variation.Name, checkMark)
			})
		}
	}

}
