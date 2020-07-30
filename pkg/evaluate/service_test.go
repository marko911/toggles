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
	e EvaluationRequest
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

	userTarget := []byte(`{
		"id": { "$oid": "5f09d08d40a5b800068a5d88" },
		"name": "Young chicks",
		"key": "hey-ladies",
		"enabled": true,
		"variations": [
			{
				"name": "On",
				"percent": 100,
				"users": ["jenny@hey.com", "mary@hey.com"]
			},
			{ "name": "Off", "percent": 0 }
		],
	
		"tenant": { "$oid": "5ef5f06a4fc7eb0006772c49" }
	}
	`)

	defaultTargetFlag := []byte(`{
		"id": { "$oid": "5f09d08d40a5b800068a5d88" },
		"name": "alpha users",
		"key": "alpha-users",
		"enabled": true,
		"variations": [
			{ "name": "Red", "percent": 100 },
			{ "name": "Blue", "percent": 0 },
			{ "name": "Purple", "percent": 0 }
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
		
		}	
	`)

	tests := []testCase{
		{
			"flag with user targeting",
			fields{r: &mock.EvaluateByte{Flag: userTarget}},
			args{EvaluationRequest{"hey-ladies", user{Key: "jenny@hey.com"}}},
			&EvaluationResult{
				models.User{Key: "jenny@hey.com"},
				&models.Variation{Name: "On", Percent: 100, UserKeys: []string{"jenny@hey.com", "mary@hey.com"}},
				bson.ObjectIdHex("5f09d08d40a5b800068a5d88"),
			},
			false,
		},
		{
			"flag default target Red variation",
			fields{r: &mock.EvaluateByte{Flag: defaultTargetFlag}},
			args{EvaluationRequest{"hey-ladies", map[string]interface{}{
				"key": "jenny@hey.com",
				"attributes": map[string]interface{}{
					"groups": []string{"ladies"},
				},
			}}},
			&EvaluationResult{
				models.User{Key: "jenny@hey.com", Attributes: map[string]interface{}{
					"groups": []string{"ladies"},
				}},
				&models.Variation{Name: "Red", Percent: 100},
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

type distCase struct {
	name       string
	fields     fields
	attributes map[string]interface{}
}

func TestMatchingDistributions(t *testing.T) {

	tests := []distCase{
		{
			"target matched and distributed rollout",
			fields{r: &mock.EvaluateByte{Flag: []byte(`{
			"id": { "$oid": "5f09d08d40a5b800068a5d88" },
			"name": "alpha users",
			"key": "alpha-users",
			"enabled": true,
			"variations": [
				{ "name": "Red", "percent": 100 },
				{ "name": "Blue", "percent": 0 },
				{ "name": "Purple", "percent": 0 }
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
			
			}`)}},
			map[string]interface{}{
				"groups": []string{"alpha users"},
			},
		},
		{
			"default variations distributed rollout",
			fields{r: &mock.EvaluateByte{Flag: []byte(`{
			"id": { "$oid": "5f09d08d40a5b800068a5d88" },
			"name": "alpha users",
			"key": "alpha-users",
			"enabled": true,
			"variations": [
				{ "name": "Red", "percent": 15 },
				{ "name": "Blue", "percent": 45 },
				{ "name": "Purple", "percent": 40 }
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
			
			}`)}},
			map[string]interface{}{
				"groups": []string{"none"},
			},
		},
	}

	t.Log("\nGiven the need to evenly distribute variations")
	for _, tt := range tests {
		t.Logf("\n\tWith %s", tt.name)

		{
			s := &service{tt.fields.r}
			randomKeyCharset := []byte("123456789abcdefghijkmnopqrstuvwxyz")
			is := is.New(t)

			var red, purple, blue int

			numUserS := 1000
			// generate random users eval requests
			for i := 1; i < numUserS; i++ {
				u := map[string]interface{}{
					"key":        uniuri.NewLenChars(uniuri.StdLen, randomKeyCharset),
					"attributes": tt.attributes,
				}

				e := EvaluationRequest{"alpha-users", u}

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
				t.Logf("\n\t\tWhen checking \"%s\" for number of variations", k)
				{
					diff := getErrorDiff(300, v)
					if diff > errorMargin {
						t.Fatalf("Distribution error margin exceeded for %s: got %v want lte %v", k, diff, errorMargin)
					}
					t.Logf("\n\t\t\tShould receive error margin less than %v percent for variation %s %v",
						errorMargin, k, checkMark)
				}
			}

		}
	}

}
