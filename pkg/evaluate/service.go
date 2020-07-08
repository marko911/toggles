package evaluate

import (
	"errors"
	"fmt"
	"toggle/server/pkg/models"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// EvaluationRequest represents an a client request on a flag key
type EvaluationRequest struct {
	FlagKey string      `json:"flagKey"`
	User    models.User `json:"user"`
}

// EvaluationData wraps User into interface for passing to evaluation
type EvaluationData struct {
	FlagKey string
	User    interface{}
}

// Service runs evaluation operations on client feature flag requests
type Service interface {
	Evaluate(e EvaluationData) (*models.EvaluationResult, error)
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
func (s *service) Evaluate(e EvaluationData) (*models.EvaluationResult, error) {
	flag, err := s.r.GetFlag(e.FlagKey)
	if err != nil {
		logrus.Error("Error Getting flag", err)
		return nil, err
	}

	var u models.User
	err = mapstructure.Decode(e.User, &u)
	if err != nil {
		return nil, errors.New("Cant cast user to User model from request data")
	}
	if v := e.VariationFromUserTargeting(flag); v != nil {
		return &models.EvaluationResult{User: u, Variation: v, FlagID: flag.ID}, nil
	}
	fmt.Println("FLAG TARGETS", flag.Targets[0].Rules)
	v, err := e.MatchFlagTarget(flag.Targets)
	if err != nil {
		return nil, err
	}
	return &models.EvaluationResult{User: u, Variation: v, FlagID: flag.ID}, nil

}

// VariationFromUserTargeting checks to see if user has been specifically targeted
// in a flag configuration and returns the varition if true
func (e *EvaluationData) VariationFromUserTargeting(f *models.Flag) *models.Variation {
	for _, variation := range f.Variations {
		if len(variation.UserKeys) > 0 {
			for _, key := range variation.UserKeys {
				user, ok := e.User.(models.User)
				if !ok {
					logrus.Error("cant cast user from request")
					return nil
				}
				if key == user.Key {
					return &variation
				}
			}
		}
	}
	return nil
}
