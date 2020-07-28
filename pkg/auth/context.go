package auth

import "context"

// FromContext returns Authorizer from context
func FromContext(c context.Context) *Authorizer {
	return c.Value(ServiceKey).(*Authorizer)
}
