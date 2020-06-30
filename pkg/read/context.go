package read

import "context"

// ServiceKey is key used for context binding
const ServiceKey string = "readCTXKey"

// FromContext returns the read service from context
func FromContext(c context.Context) Service {
	return c.Value(ServiceKey).(Service)
}
