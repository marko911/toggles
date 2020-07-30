package message

import "context"

// ServiceKey is key used for context binding
const ServiceKey string = "messageCTXKEy"

// FromContext returns the messaging service from context
func FromContext(c context.Context) Service {
	return c.Value(ServiceKey).(Service)
}
