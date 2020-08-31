.DEFAULT_GOAL := help
.PHONY: help

project?=featflags-369-superseed
tag?=gcr.io/$(project)/flags-server:dev
args?=

test:
	go test ./...

build:
	docker build -t $(tag) \
	--build-arg build_time=$(shell date -u '+%Y-%m-%d_%H:%M:%S') \
	--build-arg git_commit=$(shell git rev-parse --short HEAD) -f ./docker/Dockerfile .

push:
	docker push $(tag)

release: test build push

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1\3/p' \
	| column -t  -s ' '

run: ## Spin up dev env via docker containers
		docker-compose -f ./docker/docker-compose.yml up --build

