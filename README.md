<!--
Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
   
    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Dell Container Storage Modules (CSM)

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](docs/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/github/license/dell/csm)](LICENSE)
[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/dell/csm?include_prereleases&label=latest&style=flat-square)](https://github.com/dell/csm/releases/latest)

Dell Container Storage Modules (CSM) is an open-source suite of Kubernetes storage enablers for Dell products.

For documentation, please visit [Container Storage Modules documentation](https://dell.github.io/csm-docs/).

## Table of Contents

* [Code of Conduct](./docs/CODE_OF_CONDUCT.md)
* [Maintainer Guide](./docs/MAINTAINER_GUIDE.md)
* [Committer Guide](./docs/COMMITTER_GUIDE.md)
* [Contributing Guide](./docs/CONTRIBUTING.md)
* [List of Adopters](./docs/ADOPTERS.md)
* [Support](./docs/SUPPORT.md)
* [Security](./docs/SECURITY.md)
* [Building](#building)
* [Container Storage Modules - Components](#container-storage-modules---components)
* [About](#about)

## Building
This project includes the base container image definition for the
[Container Storage Modules - Components](#container-storage-modules---components).

To build the image, some requirements must be met:
* The supported build environment is restricted to RedHat Enterprise Linux version 9.0 and above
* buildah is used to build the container and must be installed

Once the requirements above are met, the image can be build via:
`make docker`

Note: Due to the way that buildah operates, you may see warnings (or errors) in the output of `make docker`. 
The following messages can be safely ignored:
* `/proc/ is not mounted. This is not a supported mode of operation. Please fix
your invocation environment to mount /proc/ and /sys/ properly. Proceeding anyway.` The /proc filesystem is not needed for the image creation.
* `[Errno 13] Permission denied: '/var/log/rhsm/rhsm.log' - Further logging output will be written to stderr`. This error comes from `dnf` as it is unable to write to the log file. Further errors from dnf will be sent to stderr instead of the log. 


## Container Storage Modules - Components
* [Dell Container Storage Modules (CSM) for Authorization](https://github.com/dell/karavi-authorization)
* [Dell Container Storage Modules (CSM) for Observability](https://github.com/dell/karavi-observability)
* [Dell Container Storage Modules (CSM) for Replication](https://github.com/dell/csm-replication)
* [Dell Container Storage Modules (CSM) for Resiliency](https://github.com/dell/karavi-resiliency)
* [CSI Driver for Dell PowerFlex](https://github.com/dell/csi-powerflex)
* [CSI Driver for Dell PowerMax](https://github.com/dell/csi-powermax)
* [CSI Driver for Dell PowerScale](https://github.com/dell/csi-powerscale)
* [CSI Driver for Dell PowerStore](https://github.com/dell/csi-powerstore)
* [CSI Driver for Dell Unity](https://github.com/dell/csi-unity)
* [COSI Driver](https://github.com/dell/cosi)
  
## About

Dell Container Storage Modules (CSM) is 100% open source and community-driven. All components are available
under [Apache 2 License](https://www.apache.org/licenses/LICENSE-2.0.html) on
GitHub.
