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

UBIBASE=""
CSMBASE=""
PACKAGES=()

function build() {
    echo "Building base image from ${UBIBASE}"
    echo "And creating ${CSMBASE}"
    echo "With packages of: ${PACKAGES[@]}"

    microcontainer=$(buildah from "${UBIBASE}")
    micromount=$(buildah mount $microcontainer)
    dnf install \
      --installroot $micromount \
      --releasever=9 \
      --nodocs \
      --setopt install_weak_deps=false \
      -y \
      --setopt=reposdir=/etc/yum.repos.d/ \
      ${PACKAGES[@]}

    dnf clean all \
      --installroot $micromount

    buildah umount $microcontainer
    buildah commit $microcontainer "${CSMBASE}"
}

# Parse command line arguments
while getopts "u:t:" opt; do
  case $opt in
    u)
      UBIBASE="$OPTARG"
      ;;
    t)
      CSMBASE="$OPTARG"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      ;;
  esac
done

# Remove the parsed options from the argument list
shift $((OPTIND-1))

# Store the remaining arguments in the UNNAMED_ARGS array
PACKAGES=("$@")

# Use the parsed values
if [ -z "$UBIBASE" ]; then
  echo "Error: UBIBASE is not set"
  exit 1
fi

if [ -z "$CSMBASE" ]; then
  echo "Error: CSMBASE is not set"
  exit 1
fi

build
