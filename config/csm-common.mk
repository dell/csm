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
# REGISTRY can be set as an environment variable (e.g., registry.example.com/path)
# Defaults to Red Hat public registry if not specified
# UBI_VERSION can be set to change UBI version (e.g., ubi9, ubi10)
# Defaults to ubi9 if not specified
REGISTRY ?= registry.access.redhat.com
UBI_VERSION ?= ubi9
UBI_BASEIMAGE=$(REGISTRY)/$(UBI_VERSION)/ubi-micro@sha256:b498b3ea26111ab4b81d65139f2ebd2ef9a2abb7a4588b7fdcc54889f95e9caa

# --- CSM_BASEIMAGE: Specifies the common baseimage that is used for all CSM images.
CSM_BASEIMAGE=quay.io/dell/container-storage-modules/csm-base-image:nightly

# --- DEFAULT_BASEIMAGE: Specifies the default image for repositories not yet using the CSM_BASEIMAGE.
# --- Repositories should switch to using the CSM_BASEIMAGE to use the new common base image.
DEFAULT_BASEIMAGE=${UBI_BASEIMAGE}

# --- DEFAULT_GOVERSION: Specifies the default version of go.
DEFAULT_GOVERSION=1.26

# --- DEFAULT_GOIMAGE: Specifies the default Image to be used for building go components in a multi-stage docker file.
DEFAULT_GOIMAGE=golang:${DEFAULT_GOVERSION}
