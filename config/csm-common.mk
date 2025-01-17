# Copyright © 2024-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#      http:#www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This file is intended to be included from within a Makefile,
# so is restricted to Makefile syntax

# Common settings for all CSM components and images.
# Update this file when the image versions change, and it will be automatically rolled out across all components.

# --- UBI_BASEIMAGE: Value of the UBI image to be used as a base for all images.
UBI_BASEIMAGE=registry.access.redhat.com/ubi9/ubi-micro@sha256:7f376b75faf8ea546f28f8529c37d24adcde33dca4103f4897ae19a43d58192b

# --- CSMBASEIMAGE settings: Define the values for a common UBI-micro based image.
# registry
CSMBASEIMAGE_REGISTRY=quay.io
# organization
CSMBASEIMAGE_NAMESPACE=dell/container-storage-modules
# image
CSMBASEIMAGE_NAME=csm-base-image
# tag
CSMBASEIMAGE_TAG_NEWEST=nightly
# full image name, without tag
CSMBASEIMAGE_IMAGE=${CSMBASEIMAGE_REGISTRY}/${CSMBASEIMAGE_NAMESPACE}/${CSMBASEIMAGE_NAME}

# Set the CSMBASEIMAGE_TAG_NEWEST if no semantically versioned tags were found.
ifeq ($(CSMBASEIMAGE_TAG_NEWEST),)
export CSMBASEIMAGE_TAG_NEWEST=nightly
endif

# --- CSM_BASEIMAGE: Specifies the common baseimage that is used for all CSM images.
CSM_BASEIMAGE=${CSMBASEIMAGE_IMAGE}:${CSMBASEIMAGE_TAG_NEWEST}

# --- DEFAULT_BASEIMAGE: Specifies the default image for repositories not yet using the CSM_BASEIMAGE.
# --- Repositories should switch to using the CSM_BASEIMAGE to use the new common base image.
DEFAULT_BASEIMAGE=${UBI_BASEIMAGE}

# --- DEFAULT_GOVERSION: Specifies the default version of go.
DEFAULT_GOVERSION=1.23

# --- DEFAULT_GOIMAGE: Specifies the default Image to be used for building go components in a multi-stage docker file.
DEFAULT_GOIMAGE=golang:${DEFAULT_GOVERSION}
