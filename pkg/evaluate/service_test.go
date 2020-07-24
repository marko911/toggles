package evaluate

import (
	"math"
	"testing"
	"toggle/server/pkg/mock"
	"toggle/server/pkg/models"

	"github.com/cheekybits/is"
	"github.com/dchest/uniuri"
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

var errorMargin float64 = 3

func TestMatchingDistributions(t *testing.T) {

	t.Log("Given the need to evenly distribute variations")
	{
		s := &service{
			r: &mock.EvaluateByte{Flag: []byte(`{
				"id": { "$oid": "5f09d08d40a5b800068a5d88" },
				"name": "alpha users",
				"key": "alpha-users",
				"enabled": true,
				"variations": [
					{ "name": "On", "percent": 0 },
					{ "name": "Off", "percent": 100 }
    		],
				"targets": [
					{
						"rules": [
							{ "attribute": "groups", "operator": "CONTAINS", "value": "\"alpha users\"" }
						],
						"variations": [
							{ "name": "Red", "percent": 30 },
							{ "name": "Blue", "percent": 50 },
							{ "name": "Purple", "percent": 20 }
						]
					}
    		]
				
				}`)},
		}
		randomKeyCharset := []byte("123456789abcdefghijkmnopqrstuvwxyz")
		is := is.New(t)

		var red, purple, blue int

		numUserS := 1000
		// generate random users eval requests
		for i := 1; i < numUserS; i++ {
			u := map[string]interface{}{
				"key": uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
				"attributes": map[string]interface{}{
					"groups": []string{"alpha users"},
				},
			}

			e := EvaluationData{"alpha-users", u}

			got, err := s.Evaluate(e)

			is.NoErr(err)

			if got.Variation.Name == "Red" {
				red++
			}

			if got.Variation.Name == "Blue" {
				blue++
			}

			if got.Variation.Name == "Purple" {
				purple++
			}
		}

		getErrorDiff := func(target, actual int) float64 {
			diff := math.Abs(float64(target - actual))
			percent := (diff / float64(numUserS)) * 10
			return math.Round(percent)
		}

		for k, v := range map[string]int{"red": red, "blue": blue, "purple": purple} {
			t.Logf("\tWhen checking \"%s\" for number of variations", k)
			{
				diff := getErrorDiff(300, v)
				if diff > errorMargin {
					t.Fatalf("Distribution error margin exceeded for %s: got %v want lte %v", k, diff, errorMargin)
				}
				t.Logf("\t\tShould receive error margin less than %v percent for variation %s %v",
					errorMargin, k, checkMark)
			}
		}

	}

}
