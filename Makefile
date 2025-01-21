SERVICE_NAME := metric-client
USER_ID ?= $(shell id -u)
GROUP_ID ?= $(shell id -g)
DOCKER_COMPOSE := DOCKER_BUILDKIT=1 USER_ID=$(USER_ID) GROUP_ID=$(GROUP_ID) $(shell which docker) compose
DOCKER := DOCKER_BUILDKIT=1 $(shell which docker)

ifeq ($(shell which -s go && echo true), true)
OS := $(shell go env GOOS)
ARCHITECTURE := $(shell go env GOARCH)
PACKAGES_PATH := $(shell go env GOMODCACHE)
else
OS := linux
ARCHITECTURE := amd64
PACKAGES_PATH := $(HOME)/go/pkg/mod
endif

.PHONY: start
start:
	$(DOCKER_COMPOSE) up --build

.PHONY: stop
stop:
	$(DOCKER_COMPOSE) stop

.PHONY: down
down:
	$(DOCKER_COMPOSE) down --remove-orphans

.PHONY: dc
dc:
	$(DOCKER_COMPOSE) $(c)

.PHONY: sh
sh:
	$(DOCKER_COMPOSE) exec app sh