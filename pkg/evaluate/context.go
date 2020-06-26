package evaluate

import "context"

// ServiceKey is key used for context binding
const ServiceKey string = "evaluateCTXKey"

// FromContext returns the create service from context
func FromContext(c context.Context) Service {
	return c.Value(ServiceKey).(Service)
}
