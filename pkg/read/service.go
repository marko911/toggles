package read

import (
	"toggle/server/pkg/models"

	"github.com/globalsign/mgo/bson"
)

// Service provides read operations
type Service interface {
	GetFlags(models.Tenant) ([]models.Flag, error)
	GetFlag(key string) (*models.Flag, error)
	GetFlagStats(id bson.ObjectId) (*models.FlagStats, error)
	GetSegments(models.Tenant) ([]models.Segment, error)
	GetUsers(models.Tenant) ([]models.User, error)
	GetTenant(key string) *models.Tenant
	GetEvals() ([]models.Evaluation, error)
	GetFlagEvals(bson.ObjectId, int, int) ([]models.Evaluation, int, error)
	GetTenantFromAPIKey(apiKey string) *models.Tenant
}

// Repository handles fetching persisted entities
type Repository interface {
	GetFlags(models.Tenant) ([]models.Flag, error)
	GetFlag(key string) (*models.Flag, error)
	GetFlagStats(id bson.ObjectId) (*models.FlagStats, error)
	GetSegments(models.Tenant) ([]models.Segment, error)
	GetUsers(models.Tenant) ([]models.User, error)
	GetTenant(key string) *models.Tenant
	GetEvals() ([]models.Evaluation, error)
	GetFlagEvals(bson.ObjectId, int, int) ([]models.Evaluation, int, error)
	GetTenantFromAPIKey(apiKey string) *models.Tenant
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

func (s *service) GetTenantFromAPIKey(apiKey string) *models.Tenant {
	return s.r.GetTenantFromAPIKey(apiKey)
}

func (s *service) GetFlag(key string) (*models.Flag, error) {
	flag, err := s.r.GetFlag(key)
	if err != nil {
		return nil, err
	}
	return flag, nil
}

func (s *service) GetEvals() ([]models.Evaluation, error) {
	evals, err := s.r.GetEvals()
	if err != nil {
		return nil, err
	}
	return evals, nil
}

func (s *service) GetFlagEvals(id bson.ObjectId, page int, limit int) ([]models.Evaluation, int, error) {
	evals, c, err := s.r.GetFlagEvals(id, page, limit)
	if err != nil {
		return nil, -1, err
	}
	return evals, c, nil
}

func (s *service) GetFlagStats(id bson.ObjectId) (*models.FlagStats, error) {
	stats, err := s.r.GetFlagStats(id)
	if err != nil {
		return nil, err
	}
	return stats, err
}
