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
DEFAULT_BASEIMAGE="registry.access.redhat.com/ubi9/ubi-micro@sha256:f6e0a71b7e0875b54ea559c2e0a6478703268a8d4b8bdcf5d911d0dae76aef51"

# --- DEFAULT_GOVERSION: Specifies the default version of go
DEFAULT_GOVERSION="1.23"

# --- DEFAULT_GOIMAGE: Specifies the default Image to be used for building go components in a multi-stage docker file
DEFAULT_GOIMAGE="golang:${DEFAULT_GOVERSION}"
