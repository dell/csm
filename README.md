<!--
Copyright (c) 2021-2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# :lock: **Important Notice**
Starting with the release of **Container Storage Modules v1.16.0**, this repository will no longer be maintained as an open source project. Future development will continue under a closed source model. This change reflects our commitment to delivering even greater value to our customers by enabling faster innovation and more deeply integrated features with the Dell storage portfolio.<br>
For existing customers using Dell’s Container Storage Modules, you will continue to receive:
* **Ongoing Support & Community Engagement**<br>
       You will continue to receive high-quality support through Dell Support and our community channels. Your experience of engaging with the Dell community remains unchanged.
* **Streamlined Deployment & Updates**<br>
        Deployment and update processes will remain consistent, ensuring a smooth and familiar experience.
* **Access to Documentation & Resources**<br>
       All documentation and related materials will remain publicly accessible, providing transparency and technical guidance.
* **Continued Access to Current Open Source Version**<br>
       The current open-source version will remain available under its existing license for those who rely on it.

Moving to a Enterprise model allows Dell’s development team to accelerate feature delivery and enhance integration across our Enterprise Kubernetes Storage solutions ultimately providing a more seamless and robust experience.<br>
We deeply appreciate the contributions of the open source community and remain committed to supporting our customers through this transition.<br>

For questions or access requests, please contact the maintainers via [Dell Support](https://www.dell.com/support/kbdoc/en-in/000188046/container-storage-interface-csi-drivers-and-container-storage-modules-csm-how-to-get-support).

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

## Support
For any issues, questions or feedback, please contact [Dell support](https://www.dell.com/support/incidents-online/en-us/contactus/product/container-storage-modules).

## About
Dell Container Storage Modules (CSM) is 100% open source and community-driven. All components are available
under [Apache 2 License](https://www.apache.org/licenses/LICENSE-2.0.html) on
GitHub.
