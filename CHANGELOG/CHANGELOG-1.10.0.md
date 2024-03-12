<!--toc-->
- [v1.10.0](#v1100)
  - [Changelog since v1.9.3](#changelog-since-v193)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
 

# v1.10.0 

## Changelog since v1.9.3 

## Known Issues 
- Resource quotas may not work properly with the CSI PowerFlex driver. PowerFlex is only able to assign storage in 8Gi chunks, so if a create volume call is made with a size not divisible by 8Gi, CSI-PowerFlex will round up to the next 8Gi boundary when it provisions storage -- however, the resource quota will not record this size but rather the original size in the create request. This means that, for example, if a 10Gi resource quota is set, and a user provisions 10 1Gi PVCs, 80Gi of storage will actually be allocated, which is well over the amount specified in the resource quota. For now, because of this, users should only provision volumes in 8Gi-divisible chunks if they want to use resource quotas.

## Changes by Kind 

### Deprecation 

### Features 

- CSM 1.10 release specific changes. ([#1091](https://github.com/dell/csm/issues/1091))
- Fixing the linting, formatting and vetting issues. ([#926](https://github.com/dell/csm/issues/926))
- Support PowerStore v3.6. ([#1129](https://github.com/dell/csm/issues/1129))
- Sample YAML for storage class creation missing in CSM Operator documentation. ([#1105](https://github.com/dell/csm/issues/1105))

### Bugs 

- Resource quota bypass. ([#1163](https://github.com/dell/csm/issues/1163))
- Operator known issue for offline bundle sidecar images should have examples for all platforms. ([#1180](https://github.com/dell/csm/issues/1180))
- PowerMax : Metro: Failed to find Remote Symm WWN. ([#1175](https://github.com/dell/csm/issues/1175))
- Kubelet Configuration Directory setting should not have a comment about default value being None. ([#1174](https://github.com/dell/csm/issues/1174))
- Documentation : Multipath related instructions are missing in Powerstore prerequisites. ([#1142](https://github.com/dell/csm/issues/1142))
- Cert-csi tests are not reporting the passed testcases in K8S E2E tests. ([#1140](https://github.com/dell/csm/issues/1140))
- PowerScale : Driver failing to re-authenticate if session cookies are expired. ([#1134](https://github.com/dell/csm/issues/1134))
- Inactive Tutorials button. ([#1121](https://github.com/dell/csm/issues/1121))
- CSI Powermax: Driver fails to restore snapshot to Metro Volumes. ([#1115](https://github.com/dell/csm/issues/1115))
- The csm-isilon-controller keeps getting panic and is restarting. ([#1104](https://github.com/dell/csm/issues/1104))
- the `nasName` parameter in the powerflex secret is now mandatory. ([#1101](https://github.com/dell/csm/issues/1101))
- Powerstore sanity tests are not working. ([#1097](https://github.com/dell/csm/issues/1097))
- CSM Operator offline install powerflex csi driver sidecar trying to pull from registry.k8s.io. ([#1094](https://github.com/dell/csm/issues/1094))
- PowerFlex driver fails to start on RKE. ([#1086](https://github.com/dell/csm/issues/1086))
- Gopowerstore - Multiple issues in mockfile. ([#1084](https://github.com/dell/csm/issues/1084))
- CSM driver repositories reference CSI Operator. ([#1081](https://github.com/dell/csm/issues/1081))
