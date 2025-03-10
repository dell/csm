# default target
all: help

# include an overrides file, which sets up default values and allows user overrides
include overrides.mk

# variables
BASE_IMAGE_PACKAGES=acl \
					gnutls \
					device-mapper-multipath \
					e2fsprogs \
					gnutls \
					gzip \
					bzip2 \
					hostname \
					kmod \
					libaio \
					libblockdev \
					libuuid \
					libxcrypt-compat \
					nettle \
					nfs-utils \
					nfs4-acl-tools \
					numactl \
					openssl \
					rpm \
					systemd \
					tar \
					util-linux \
					which \
					xfsprogs

# Help target, prints useful information
help:
	@echo
	@echo "The following targets are commonly used:"
	@echo
	@echo "docker           - Builds the container image"
	@echo

# Clean the build
clean:
	rm -f core/core_generated.go
	rm -f semver.mk
	go clean

# Generates the docker container (but does not push)
docker:
	$(eval include config/csm-common.mk)
	$(eval include semver.mk)
	@echo "Building base image from $(UBI_BASEIMAGE) and loading dependencies..."
	cd base-image && ./build-base-image.sh -u $(UBI_BASEIMAGE) -t $(REGISTRY)/$(IMAGENAME):$(IMAGETAG) $(BASE_IMAGE_PACKAGES)
	@echo "Built base image: $(REGISTRY)/$(IMAGENAME):$(IMAGETAG)"
