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
# limitations under the License.

# --- DOCKERFILE ARGUMENTS

# reference to the UBI micro image
ARG UBIMICRO=unknown

# ubi image corresponding to the micro version
ARG UBIBUILDER=unknown

# mount dir for installing option packages, should not need to be changed
ARG MICRODIR=/microdir

# list of packages to install
ARG PACKAGES_TO_INSTALL="rpm"



# --- microbase: image provided by RedHat as the UBI micro base image
FROM $UBIMICRO AS microbase

# --- microbuilder: Used to add packages to the microbase image since no package manager exists in micro
FROM $UBIBUILDER AS microbuilder
ARG MICRODIR
ARG PACKAGES_TO_INSTALL
RUN mkdir ${MICRODIR}
COPY --from=microbase / ${MICRODIR}
RUN yum install \
  --installroot ${MICRODIR} \
  --releasever 9 \
  --setop install_weak_deps=false \
  -y \
  ${PACKAGES_TO_INSTALL}
RUN yum clean all \
  --installroot ${MICRODIR}

# --- final: image, build from microabse with packages installed by microbuilder
FROM $UBIMICRO AS final
ARG MICRODIR
COPY --from=microbuilder ${MICRODIR}/ .
