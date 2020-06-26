# Feature Flags API


An API for creating feature flags and user segments, deployable on Kubernetes.

Uses NATS for connecting potential breakout services + managed MongoDB. 

Implements a hexagonal architecture with most service dependencies injected in `server.go` 