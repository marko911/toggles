package evaluate

import "toggle/server/pkg/models"

// Service runs evaluation operations on client feature flag requests
type Service interface {
	Evaluate(e models.Evaluation) (*models.EvaluationResult, error)
}

// Repository holds all persisted data related to flags, users, attributes, segments
type Repository interface {
	GetFlag(key string) (*models.Flag, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Evaluate(e models.Evaluation) (*models.EvaluationResult, error) {
	flag, err := s.r.GetFlag(e.FlagKey)
	if err != nil {
		return nil, err
	}

	if eval := matchesUserTargeting(e.User, flag); eval != nil {
		res := eval.performEvaluation()
		return res
	}

}
