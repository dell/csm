# overrides file
# this file, included from the Makefile, will overlay default values with environment variables
#

# DEFAULT values
DEFAULT_REGISTRY="sample_registry"
DEFAULT_IMAGENAME="csm-base-image"
DEFAULT_BUILDSTAGE="final"
DEFAULT_IMAGETAG="test"

# set the GOIMAGE if needed
ifeq ($(GOIMAGE),)
export GOIMAGE="$(DEFAULT_GOIMAGE)"
endif

# set the REGISTRY if needed
ifeq ($(REGISTRY),)
export REGISTRY="$(DEFAULT_REGISTRY)"
endif

# set the IMAGENAME if needed
ifeq ($(IMAGENAME),)
export IMAGENAME="$(DEFAULT_IMAGENAME)"
endif

#set the IMAGETAG if needed
ifeq ($(IMAGETAG),) 
export IMAGETAG="$(DEFAULT_IMAGETAG)"
endif

# set the BUILDSTAGE if needed
ifeq ($(BUILDSTAGE),)
export BUILDSTAGE="$(DEFAULT_BUILDSTAGE)"
endif

# figure out if podman or docker should be used (use podman if found)
ifneq (, $(shell which podman 2>/dev/null))
export BUILDER=podman
else
export BUILDER=docker
endif

# target to print some help regarding these overrides and how to use them
overrides-help:
	@echo
	@echo "The following environment variables can be set to control the build"
	@echo
	@echo "GOIMAGE   - The version of Go to build with, default is: $(DEFAULT_GOIMAGE)"
	@echo "              Current setting is: $(GOIMAGE)"
	@echo "REGISTRY    - The registry to push images to, default is: $(DEFAULT_REGISTRY)"
	@echo "              Current setting is: $(REGISTRY)"
	@echo "IMAGENAME   - The image name to be built, defaut is: $(DEFAULT_IMAGENAME)"
	@echo "              Current setting is: $(IMAGENAME)"
	@echo "IMAGETAG    - The image tag to be built, default is an empty string which will determine the tag by examining annotated tags in the repo."
	@echo "              Current setting is: $(IMAGETAG)"
	@echo "BUILDSTAGE  - The Dockerfile build stage to execute, default is: $(DEFAULT_BUILDSTAGE)"
	@echo "              Stages can be found by looking at the Dockerfile"
	@echo "              Current setting is: $(BUILDSTAGE)"
	@echo
        
	

