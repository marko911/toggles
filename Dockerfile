FROM golang:1.13.8-alpine3.11 AS builder

ARG git_commit
ARG build_time

WORKDIR /toggle/server

# Copy all source files over.
COPY . .
# Install dependencies
RUN go mod download

RUN COMMIT=${git_commit} BUILD_TIME=${build_time} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags "-s -w -X toggle/server/version.Commit=${git_commit} -X toggle/server/version.BuildTime=${build_time}" \
  toggle/server/cmd/toggle

FROM alpine:latest

# Copy over the source and compiled binary.
WORKDIR /toggle

COPY --from=builder /toggle/server /toggle/
EXPOSE 8080

ENTRYPOINT ./toggle
