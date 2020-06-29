package create

import "toggle/server/pkg/models"

type Payload []models.Flag

type Event int

const (
	Done Event = iota

	FlagKeyExists

	Failed
)

// Service provides create operations
type Service interface {
	CreateFlag(*models.Flag) error
	CreateSegment(*models.Segment) error
	CreateUser(*models.User) error
}

type Repository interface {
	InsertFlag(*models.Flag) error
	InsertSegment(*models.Segment) error
	InsertUser(*models.User) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateFlag(f *models.Flag) error {
	err := s.r.InsertFlag(f)
	return err
}

func (s *service) CreateSegment(seg *models.Segment) error {
	err := s.r.InsertSegment(seg)
	return err
}

func (s *service) CreateUser(u *models.User) error {
	err := s.r.InsertUser(u)
	return err
}
