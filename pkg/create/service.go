package create

import "toggle/server/pkg/models"

// Service provides create operations
type Service interface {
	CreateFlag(*models.Flag) error
	CreateSegment(*models.Segment) error
	CreateUser(*models.User) error
}

// Repository handles persistance of entity data
type Repository interface {
	InsertFlag(*models.Flag) error
	InsertSegment(*models.Segment) error
	InsertUser(*models.User) error
}

type service struct {
	r Repository // outbound port
}

// NewService returns a creation service
func NewService(r Repository) Service {
	return &service{r}
}

// CreateFlag creates a new flag in repository
func (s *service) CreateFlag(f *models.Flag) error {
	err := s.r.InsertFlag(f)
	return err
}

// CreateSegment creates a new segment in repository
func (s *service) CreateSegment(seg *models.Segment) error {
	err := s.r.InsertSegment(seg)
	return err
}

// CreateUser creates a new user in repository
func (s *service) CreateUser(u *models.User) error {
	err := s.r.InsertUser(u)
	return err
}
