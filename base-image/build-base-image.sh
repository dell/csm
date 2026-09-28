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

# REGISTRY specifies the container registry location (optional)
# Defaults to Red Hat public registry if not specified
REGISTRY="registry.access.redhat.com"
# UBI_VERSION specifies the UBI version (e.g., ubi9, ubi10)
# Defaults to ubi9 if not specified
UBI_VERSION="ubi9"
# UBIBASE refers to the UBI-micro base image reference
# If not specified, will be constructed from REGISTRY and UBI_VERSION
UBIBASE=""
# CSMBASE specifies the fule name of the target image to build
CSMBASE=""
# PACKAGES is an array of packages that need to be installed in the built image
PACKAGES=""


# --- help: displays a short help message
function help() {
  echo "${0}: Builds an image based on a RedHat UBI Micro image, with additional packages added"
  echo "  Required Arguments:"
  echo "    -t: Name of the target image to build"
  echo "    List of package names to install"
  echo "  Optional Arguments"
  echo "    -r: Container registry location including path (e.g., registry.example.com/path)"
  echo "         Defaults to registry.access.redhat.com"
  echo "         UBIBASE is automatically constructed from REGISTRY if -u not provided"
  echo "    -v: UBI version (e.g., ubi9, ubi10)"
  echo "         Defaults to ubi9"
  echo "         UBIBASE is automatically constructed from UBI_VERSION if -u not provided"
  echo "    -u: Reference to the UBI-micro image to use as a base (optional)"
  echo "         If not provided, constructed from REGISTRY and UBI_VERSION"
  echo "    -h: Help, displays this message"
  echo ""
  echo ""
  echo "For example to build an image by adding curl and wget, invoke as:"
  echo "${0} \\"
  echo "  -t localhost/myorganization/myimage:mytag \\"
  echo "  curl wget"
  echo ""
  echo "Or with custom registry (UBIBASE auto-constructed):"
  echo "${0} \\"
  echo "  -r registry.example.com/path \\"
  echo "  -t localhost/myorganization/myimage:mytag \\"
  echo "  curl wget"
  echo
}

# --- build: Builds a container image using buildah
function build() {
  echo "Building base image from ${UBIBASE}"
  echo "And creating ${CSMBASE}"
  echo "With packages of: ${PACKAGES}"

  # export the settings
  export REGISTRY="${REGISTRY}"
  export UBIBASE="${UBIBASE}"
  export CSMBASE="${CSMBASE}"
  export PACKAGES="${PACKAGES}"

  if [ $(id -u) -eq 0 ]; then
    # if running as root, just run the script
    ./buildah-script.sh
  else
    # otherwise run in an unshared environment
    # and run the build script
    buildah unshare ./buildah-script.sh
  fi

}

# check to see if the host is RedHat Enterprise Linux as it is required
if [ ! -f /etc/redhat-release ]; then
  echo "This does not appear to be a RedHat Enterprise Linux system"
  echo "No file at /etc/redhat-release was found"
  exit 1
fi

# Parse command line arguments
while getopts "hr:v:u:t:" opt; do
  case $opt in
    r)
      REGISTRY="$OPTARG"
      ;;
    v)
      UBI_VERSION="$OPTARG"
      ;;
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

# Construct UBIBASE from REGISTRY and UBI_VERSION if not explicitly provided
if [ -z "$UBIBASE" ]; then
  UBIBASE="${REGISTRY}/${UBI_VERSION}/ubi-micro@sha256:9dbba858e5c8821fbe1a36c376ba23b83ba00f100126f2073baa32df2c8e183a"
fi

if [ -z "$CSMBASE" ]; then
  echo "Error: CSMBASE is not set"
  exit 1
fi

build
