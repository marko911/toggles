package evaluate

import (
	"toggle/server/pkg/errors"
	"toggle/server/pkg/models"

	"github.com/mitchellh/mapstructure"
)

// EvaluationRequest wraps User into interface for passing to evaluation library
type EvaluationRequest struct {
	User    interface{}
	FlagKey string `json:"flagKey"`
}

// Service runs evaluation operations on client feature flag requests
type Service interface {
	Evaluate(e EvaluationRequest, flag *models.Flag) (*models.Evaluation, error)
	MatchDefault(e EvaluationRequest, flag *models.Flag) (*models.Evaluation, error)
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
func (s *service) Evaluate(e EvaluationRequest, flag *models.Flag) (*models.Evaluation, error) {
	var u models.User
	err := mapstructure.Decode(e.User, &u)
	if err != nil {
		return nil, errors.ErrCantCastUser
	}

	if v := e.VariationFromUserTargeting(flag, &u); v != nil {
		return &models.Evaluation{Variation: v, Flag: *flag}, nil
	}

	if v, err := e.MatchFlagTarget(flag); v != nil {
		if err != nil {
			return nil, err
		}
		return &models.Evaluation{Variation: v, Flag: *flag}, nil
	}

	if v, err := e.MatchDefaultVariations(flag); v != nil {
		if err != nil {
			return nil, err
		}
		return &models.Evaluation{Variation: v, Flag: *flag}, nil
	}

	return nil, errors.ErrVariationNotFound

}

// MatchDefault matches default variation, copying from above until later
func (s *service) MatchDefault(e EvaluationRequest, flag *models.Flag) (*models.Evaluation, error) {

	if v, err := e.MatchDefaultVariations(flag); v != nil {
		if err != nil {
			return nil, err
		}
		return &models.Evaluation{Variation: v, Flag: *flag}, nil
	}

	return nil, errors.ErrVariationNotFound
}

// VariationFromUserTargeting checks to see if user has been specifically targeted
// in a flag configuration and returns the varition if true
func (e *EvaluationRequest) VariationFromUserTargeting(f *models.Flag, u *models.User) *models.Variation {
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
