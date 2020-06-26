package create

import "context"

// ServiceKey is key used for context binding
const ServiceKey string = "createCTXKey"

// FromContext returns the create service from context
func FromContext(c context.Context) Service {
	return c.Value(ServiceKey).(Service)
}
