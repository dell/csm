# default target
all: help

# include an overrides file, which sets up default values and allows user overrides
include overrides.mk

# variables
ACT_OPTIONS=--secret GITHUB_TOKEN=$(GITHUB_TOKEN) --no-cache-server --platform ubuntu-latest=ghcr.io/catthehacker/ubuntu:act-latest --github-instance github.com

# Help target, prints usefule information
help:
	@echo
	@echo "The following targets are commonly used:"
	@echo
	@echo "action-help      - Displays instructions on how to run a single github workflow locally"
	@echo "actions          - Run all workflows locally, requires https://github.com/nektos/act"
	@echo "docker           - Builds the code within a golang container and then creates the driver image"
	@echo

# Clean the build
clean:
	rm -f core/core_generated.go
	rm -f semver.mk
	go clean

# Dependencies
dependencies:
	go run core/semver/semver.go -f mk >semver.mk

# Generates container via a github workflow
docker-action:
	act pull_request \
		$(ACT_OPTIONS) \
		--job build-base-image

# Generates the docker container (but does not push)
docker: dependencies
	$(eval include config/csm-common.mk)
	$(eval include semver.mk)
	@echo "Base Images is set to: $(BASEIMAGE)"
	@echo "Building: $(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
	cd base-image && \
		$(BUILDER) build \
		-t "$(REGISTRY)/$(IMAGENAME):$(IMAGEVER)" \
		--target $(BUILDSTAGE) \
		--build-arg UBIMICRO=$(DEFAULT_BASEIMAGE) \
		--build-arg UBIBUILDER=$(DEFAULT_UBIBUILDER) \
		.

.PHONY: actions
actions: ## Run all the github action checks that run on a pull_request creation
	act -l | grep -v ^Stage | grep pull_request | awk '{print $$2}' | while read WF; do act pull_request $(ACT_OPTIONS) --job "$${WF}"; done

.PHONY: action-help
action-help: ## Echo instructions to run one specific workflow locally
	@echo "GitHub Workflows can be run locally with the following command:"
	@echo "act pull_request $(ACT_OPTIONS) --job <jobid>"
	@echo
	@echo "Where '<jobid>' is a Job ID returned by the command:"
	@echo "act -l"
	@echo
	@echo "NOTE: if act if not installed, it can be from https://github.com/nektos/act"

