#!/bin/bash
#
# Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#      http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

# This is build script that can be run in a `buildah unshare` environment

# The following environment variables are expected to be found
# UBIBASE refers to the UBI-micro base image reference
# CSMBASE specifies the fule name of the target image to build
# PACKAGES is an array of packages that need to be installed in the built image

# build the image
microcontainer=$(buildah from "${UBIBASE}")
micromount=$(buildah mount $microcontainer)
dnf install \
  --installroot $micromount \
  --releasever=9 \
  --nodocs \
  --setopt install_weak_deps=false \
  --assumeyes \
  --setopt=reposdir=/etc/yum.repos.d/ \
  ${PACKAGES}

dnf clean all \
  --installroot $micromount

buildah umount $microcontainer
buildah commit $microcontainer "${CSMBASE}"
