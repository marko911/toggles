package evaluate

import "toggle/server/pkg/models"

// EvaluationRequest represents an a client request on a flag key
type EvaluationRequest struct {
	FlagKey string      `json:"flagKey"`
	User    models.User `json:"user"`
}

// Service runs evaluation operations on client feature flag requests
type Service interface {
	Evaluate(e EvaluationRequest) (*models.EvaluationResult, error)
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
func (s *service) Evaluate(e EvaluationRequest) (*models.EvaluationResult, error) {
	flag, err := s.r.GetFlag(e.FlagKey)

	if err != nil {
		return nil, err
	}

	if v := e.VariationFromUserTargeting(flag); v != nil {
		// return variation
		return &models.EvaluationResult{User: e.User, Variation: v, FlagID: flag.ID}, nil
	}

	if v := e.MatchFlagTarget(flag.Targets); v != nil {

	}

	return nil, nil
}

// VariationFromUserTargeting checks to see if user has been specifically targeted
// in a flag configuration and returns the varition if true
func (e *EvaluationRequest) VariationFromUserTargeting(f *models.Flag) *models.Variation {
	for _, variation := range f.Variations {
		if len(variation.UserKeys) > 0 {
			for _, key := range variation.UserKeys {
				if key == e.User.Key {
					return &variation
				}
			}
		}
	}
	return nil
}
