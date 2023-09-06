<!--toc-->
- [v1.8.0](#v180)
  - [Changelog since v1.7.1](#changelog-since-v171)
  - [Known Issues](#known-issues)
    - [CSI PowerStore: Not able to create ephemeral pods in OpenShift 4.13](#csi-powerstore-not-able-to-create-ephemeral-pods-in-openshift-4.13)
    - [CSI PowerStore: In OpenShift 4.13, root user is not allowed to perform write operations on NFS shares when root squashing is enabled](#csi-powerstore-in-openshift-4.13-root-user-is-not-allowed-to-perform-write-operations-on-nfs-shares-when-root-squashing-is-enabled)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
 

# v1.8.0 

## Changelog since v1.7.1 

## Known Issues

### CSI PowerStore: Not able to create ephemeral pods in OpenShift 4.13

Ephemeral pod is not being created in OpenShift 4.13 and is failing with the error "error when creating pod: the pod uses an inline volume provided by CSIDriver csi-powerstore.dellemc.com, and the namespace has a pod security enforcement level that is lower than privileged."

This issue occurs because OpenShift 4.13 introduced the CSI Volume Admission plugin to restrict the use of a CSI driver capable of provisioning CSI ephemeral volumes during pod admission (https://docs.openshift.com/container-platform/4.13/storage/container_storage_interface/ephemeral-storage-csi-inline.html). Therefore, an additional label "security.openshift.io/csi-ephemeral-volume-profile" needs to be added to the CSIDriver object to support inline ephemeral volumes.

### CSI PowerStore: In OpenShift 4.13, root user is not allowed to perform write operations on NFS shares when root squashing is enabled

In OpenShift 4.13, the root user is not allowed to perform write operations on NFS shares, when root squashing is enabled.

The workaround for this issue is to disable root squashing by setting allowRoot: "true" in the NFS storage class.

## Changes by Kind 

### Deprecation 

### Features 

- K8S 1.28 support in CSM 1.8. ([#947](https://github.com/dell/csm/issues/947))
- Add support for Offline Install of CSM Operator in non OLM environment. ([#939](https://github.com/dell/csm/issues/939))
- Use ubi9 micro as base image. ([#922](https://github.com/dell/csm/issues/922))
- Enhancing Unity XT driver to handle API requests after the sessionIdleTimeOut in STIG mode. ([#891](https://github.com/dell/csm/issues/891))
- Make standalone helm chart available from helm repository : https://dell.github.io/dell/helm-charts. ([#877](https://github.com/dell/csm/issues/877))
- CSI 1.5 spec support -StorageCapacityTracking. ([#876](https://github.com/dell/csm/issues/876))
- CSM for PowerMax file support. ([#861](https://github.com/dell/csm/issues/861))
- CSI-PowerFlex 4.0 NFS support. ([#763](https://github.com/dell/csm/issues/763))
- CSM support for Openshift 4.13. ([#724](https://github.com/dell/csm/issues/724))
- Google Anthos 1.15 support  for PowerMax. ([#937](https://github.com/dell/csm/issues/937))
- Enhance GoPowerScale to support PowerScale Terraform Provider. ([#888](https://github.com/dell/csm/issues/888))
- Configurable Volume Attributes use recommended naming convention <prefix>/<name>. ([#879](https://github.com/dell/csm/issues/879))
- CSI 1.5 spec support: Implement Volume Limits. ([#878](https://github.com/dell/csm/issues/878))
- Use ubi-micro as base image. ([#790](https://github.com/dell/csm/issues/790))

### Bugs 

- Remove refs to deprecated io/ioutil. ([#916](https://github.com/dell/csm/issues/916))
- Unity XT: Volume Mount Hangs. ([#901](https://github.com/dell/csm/issues/901))
- Powerscale CSI driver RO PVC-from-snapshot wrong zone. ([#487](https://github.com/dell/csm/issues/487))
- CSM Operator manual installation missing volumeSnapshot CRDs as a prerequisite. ([#944](https://github.com/dell/csm/issues/944))
- cert-csi help message uses wrong name of "csi-cert". ([#938](https://github.com/dell/csm/issues/938))
- volume-group-snapshot test observes a panic when using "--namespace" parameter. ([#931](https://github.com/dell/csm/issues/931))
- "--attr" of ephemeral-volume performance test doesn't support properties file. ([#930](https://github.com/dell/csm/issues/930))
- PowerStore Replication - Delete RG request hangs. ([#928](https://github.com/dell/csm/issues/928))
- VolumeHealthMetricSuite test failure. ([#923](https://github.com/dell/csm/issues/923))
- Generating report from multiple databases and test runs failure. ([#921](https://github.com/dell/csm/issues/921))
- Update Cert-csi documentation for driver certification. ([#914](https://github.com/dell/csm/issues/914))
- Cert Manager should display tooltip about the pre-requisite.. ([#907](https://github.com/dell/csm/issues/907))
- Space is not reflecting right on Unity. ([#902](https://github.com/dell/csm/issues/902))
- Unable to pull podmon image from local repository for offline install. ([#898](https://github.com/dell/csm/issues/898))
- Documentation - Authorization. ([#895](https://github.com/dell/csm/issues/895))
- CSI driver does not verify iSCSI initiators on the array correctly. ([#849](https://github.com/dell/csm/issues/849))
- Common section for Volume Snapshot Requirements. ([#811](https://github.com/dell/csm/issues/811))
