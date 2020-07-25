package evaluate

import (
	"toggle/server/pkg/errors"
	"toggle/server/pkg/models"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// EvaluationRequest represents an a client request on a flag key
type EvaluationRequest struct {
	FlagKey string      `json:"flagKey"`
	User    models.User `json:"user"`
}

// EvaluationData wraps User into interface for passing to evaluation library
type EvaluationData struct {
	FlagKey string
	User    interface{}
}

// EvaluationResult determines what variation user is shown, can be simple true or false
// or a specific variation of flag being evaluated
type EvaluationResult struct {
	User      models.User
	Variation *models.Variation
	FlagID    bson.ObjectId `json:"flagId"`
}

// Service runs evaluation operations on client feature flag requests
type Service interface {
	Evaluate(e EvaluationData) (*EvaluationResult, error)
}

// Repository holds all persisted data related to flags, users, attributes, segments
type Repository interface {
	GetFlag(key string) (*models.Flag, error)
}

type service struct {
	r Repository
}

// NewService constructs an evaluation service
func NewService(r Repository) Service {
	return &service{r}
}

// Evaluate processes a client request and returns the variation to show to user
func (s *service) Evaluate(e EvaluationData) (*EvaluationResult, error) {

	flag, err := s.r.GetFlag(e.FlagKey)
	if err != nil {
		logrus.Error("Could not get flag with key: ", e.FlagKey)
		return nil, err
	}

	var u models.User
	err = mapstructure.Decode(e.User, &u)
	if err != nil {
		return nil, errors.ErrCantCastUser
	}

	if v := e.VariationFromUserTargeting(flag, &u); v != nil {
		return &EvaluationResult{User: u, Variation: v, FlagID: flag.ID}, nil
	}

	if v, err := e.MatchFlagTarget(flag.Targets); v != nil {
		if err != nil {
			return nil, err
		}
		return &EvaluationResult{User: u, Variation: v, FlagID: flag.ID}, nil
	}

	if v, err := e.MatchDefaultVariations(flag); v != nil {
		if err != nil {
			return nil, err
		}
		return &EvaluationResult{User: u, Variation: v, FlagID: flag.ID}, nil
	}

	return nil, errors.ErrVariationNotFound

}

// VariationFromUserTargeting checks to see if user has been specifically targeted
// in a flag configuration and returns the varition if true
func (e *EvaluationData) VariationFromUserTargeting(f *models.Flag, u *models.User) *models.Variation {
	for _, variation := range f.Variations {
		if len(variation.UserKeys) > 0 {
			for _, key := range variation.UserKeys {
				if key == u.Key {
					return &variation
				}
			}
		}
	}
	return nil
}
