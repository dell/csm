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

# --- DEFAULT_BASEIMAGE: Specifies the UBI-micro image that is used as the base for the CSM images
DEFAULT_BASEIMAGE="registry.access.redhat.com/ubi9/ubi-micro@sha256:9dbba858e5c8821fbe1a36c376ba23b83ba00f100126f2073baa32df2c8e183a"

# --- DEFAULT_UBIBUILDER: Specifies the UBI (non-micro) image that is used to install packages into the CSM base image
DEFAULT_UBIBUILDER="registry.access.redhat.com/ubi9/ubi@sha256:1ee4d8c50d14d9c9e9229d9a039d793fcbc9aa803806d194c957a397cf1d2b17"

# --- DEFAULT_GOVERSION: Specifies the default version of go
DEFAULT_GOVERSION="1.23"

# --- DEFAULT_GOIMAGE: Specifies the default Image to be used for building go components in a multi-stage docker file
DEFAULT_GOIMAGE="golang:${DEFAULT_GOVERSION}"
