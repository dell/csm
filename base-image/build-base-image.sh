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

# --- variables that are passed in via command line options

# UBIBASE refers to the UBI-micro base image reference
UBIBASE=""
# CSMBASE specifies the fule name of the target image to build
CSMBASE=""
# PACKAGES is an array of packages that need to be installed in the built image
PACKAGES=""


# --- help: displays a short help message
function help() {
  echo "${0}: Builds an image based on a RedHat UBI Micro image, with additional packages added"
  echo "  Required Arguments:"
  echo "    -u: Reference to the UBI-micro image to use as a base"
  echo "    -t: Name of the target image to build"
  echo "    List of package names to install"
  echo "  Optional Arguments"
  echo "    -h: Help, displays this message"
  echo ""
  echo ""
  echo "For example to build an image by adding curl and wget, invoke as:"
  echo "${0} \\"
  echo "  -u registry.access.redhat.com/ubi9/ubi-micro@sha256:9dbba858e5c8821fbe1a36c376ba23b83ba00f100126f2073baa32df2c8e183a \\"
  echo "  -t localhost\myorganization\myimage:mytag \\"
  echo "  curl wget"
  echo
}

# --- build: Builds a container image using buildah
function build() {
  echo "Building base image from ${UBIBASE}"
  echo "And creating ${CSMBASE}"
  echo "With packages of: ${PACKAGES}"

  # export the settings
  export UBIBASE="${UBIBASE}"
  export CSMBASE="${CSMBASE}"
  export PACKAGES="${PACKAGES}"
  # and run the build script
  buildah unshare ./buildah-script.sh
}

# check to see if the host is RedHat Enterprise Linux as it is required
if [ ! -f /etc/redhat-release ]; then
  echo "This does not appear to eb a RedHat Enterprise Linux system"
  echo "No file at /etc/redhat-release was found"
  exit 1
fi

# Parse command line arguments
while getopts "hu:t:" opt; do
  case $opt in
    u)
      UBIBASE="$OPTARG"
      ;;
    t)
      CSMBASE="$OPTARG"
      ;;
    h)
      help
      exit 0
      ;;
    \?)
      echo ""
      help
      exit 1
      ;;
  esac
done

# Remove the parsed options from the argument list
shift $((OPTIND-1))

# Store the remaining arguments in the UNNAMED_ARGS array
PACKAGES="$*"

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
