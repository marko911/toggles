package read

import (
	"toggle/server/pkg/models"
)

// Service provides read operations
type Service interface {
	GetFlags(models.Tenant) ([]models.Flag, error)
	GetFlag(key string) (*models.Flag, error)
	GetSegments(models.Tenant) ([]models.Segment, error)
	GetUsers(models.Tenant) ([]models.User, error)
	GetTenant(key string) *models.Tenant
}

// Repository handles fetching persisted entities
type Repository interface {
	GetFlags(models.Tenant) ([]models.Flag, error)
	GetFlag(key string) (*models.Flag, error)
	GetSegments(models.Tenant) ([]models.Segment, error)
	GetUsers(models.Tenant) ([]models.User, error)
	GetTenant(key string) *models.Tenant
}

type service struct {
	r Repository
}

// NewService returns a read service
func NewService(r Repository) Service {
	return &service{r}
}

// GetFlags fetches flags from repository
func (s *service) GetFlags(t models.Tenant) ([]models.Flag, error) {
	flags, err := s.r.GetFlags(t)
	if err != nil {
		return nil, err
	}
	return flags, nil
}

// GetSegments fetches segments from repository
func (s *service) GetSegments(t models.Tenant) ([]models.Segment, error) {
	segs, err := s.r.GetSegments(t)
	if err != nil {
		return nil, err
	}
	return segs, nil
}

// GetUsers fetches users from repository
func (s *service) GetUsers(t models.Tenant) ([]models.User, error) {
	users, err := s.r.GetUsers(t)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) GetTenant(key string) *models.Tenant {
	tenant := s.r.GetTenant(key)
	return tenant
}

func (s *service) GetFlag(key string) (*models.Flag, error) {
	flag, err := s.r.GetFlag(key)
	if err != nil {
		return nil, err
	}
	return flag, nil
}
