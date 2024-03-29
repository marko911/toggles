package create

import (
	"toggle/server/pkg/keygen"
	"toggle/server/pkg/models"
)

// Service provides create operations
type Service interface {
	CreateFlag(*models.Flag) error
	UpdateFlag(*models.Flag) error
	CreateSegment(*models.Segment) error
	UpdateSegment(*models.Segment) error
	CreateUser(*models.User) error
	CreateAttributes(*models.User) error
	CreateTenant(key string) (*models.Tenant, error)
	CreateEvaluation(*models.Evaluation) error
}

// Repository handles persistance of entity data
type Repository interface {
	InsertFlag(*models.Flag) error
	UpdateFlag(*models.Flag) error
	InsertSegment(*models.Segment) error
	UpdateSegment(*models.Segment) error
	InsertUser(*models.User) error
	InsertAttributes([]models.Attribute) error
	InsertTenant(*models.Tenant) (*models.Tenant, error)
	InsertEvaluation(*models.Evaluation) error
}

type service struct {
	r Repository // outbound port
}

// NewService returns a creation service
func NewService(r Repository) Service {
	return &service{r}
}

// CreateFlag creates a new flag in repository
func (s *service) CreateFlag(flag *models.Flag) error {
	err := s.r.InsertFlag(flag)
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

func (s *service) CreateAttributes(u *models.User) error {
	var attrs []models.Attribute
	for attribute := range u.Attributes {
		a := models.Attribute{
			Name: attribute,
			// Tenant: u.Tenant, //TODO: fix this
		}
		attrs = append(attrs, a)
	}

	if len(attrs) > 0 {
		err := s.r.InsertAttributes(attrs)
		return err
	}
	return nil
}

func (s *service) CreateTenant(key string) (*models.Tenant, error) {
	t := &models.Tenant{Key: key, APIKEY: keygen.GenerateToken()}

	tenant, err := s.r.InsertTenant(t)
	return tenant, err
}

func (s *service) CreateEvaluation(e *models.Evaluation) error {
	return s.r.InsertEvaluation(e)
}

func (s *service) UpdateFlag(f *models.Flag) error {
	return s.r.UpdateFlag(f)
}

func (s *service) UpdateSegment(seg *models.Segment) error {
	return s.r.UpdateSegment(seg)
}
