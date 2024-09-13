# Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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

# Common settings for all CSM components amd images.
# Update this file with the image versions change, and it will be automaticall rolled out across all components.

# --- CSMBASEIMAGE settings: Define the values for a common UBI-micro based image 
CSMBASEIMAGE_REGISTRY="docker.io"
CSMBASEIMAGE_NAMESPACE="dellemc"
CSMBASEIMAGE_NAME="csm-base-image"
CSMBASEIMAGE_TAG_NEWEST=$(curl -s "https://hub.docker.com/v2/namespaces/${CSMBASEIMAGE_NAMESPACE}/repositories/${CSMBASEIMAGE_NAME}/tags?page_size=1000" | grep -o '"name": *"[^"]*' | grep -o '[^"]*$' | grep -E '^([0-9]|v[0-9])' | sort -V | tail -n 1)
CSMBASEIMAGE_IMAGE="${CSMBASEIMAGE_REGISTRY}/${CSMBASEIMAGE_NAMESPACE}/${CSMBASEIMAGE_NAME}"

# set the CSMBASEIMAGE_TAG_NEWEST if no semanticly versioned tags were found
ifeq ($(CSMBASEIMAGE_TAG_NEWEST),)
export CSMBASEIMAGE_TAG_NEWEST="nightly"
endif

# UBI_BASEIMAGE settings: Define values for the UBI-micro image to be used as a base
UBI_BASEIMAGE="registry.access.redhat.com/ubi9/ubi-micro@sha256:9dbba858e5c8821fbe1a36c376ba23b83ba00f100126f2073baa32df2c8e183a"

# --- DEFAULT_BASEIMAGE: Specifies the UBI-micro image that is used as the base for the CSM images
DEFAULT_BASEIMAGE="${CSMBASEIMAGE_IMAGE}:${CSMBASEIMAGE_TAG_NEWEST}"

# --- DEFAULT_GOVERSION: Specifies the default version of go
DEFAULT_GOVERSION="1.23"

# --- DEFAULT_GOIMAGE: Specifies the default Image to be used for building go components in a multi-stage docker file
DEFAULT_GOIMAGE="golang:${DEFAULT_GOVERSION}"
