<!--toc-->

- [v1.13.2](#v1132)
  - [Changelog since v1.13.1](#changelog-since-v1131)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)

- [v1.13.1](#v1131)
  - [Changelog since v1.13.0](#changelog-since-v1130)
  - [Known Issues](#known-issues-1)
  - [Changes by Kind](#changes-by-kind-1)
    - [Deprecation](#deprecation-1)
    - [Features](#features-1)
    - [Bugs](#bugs-1)


- [v1.13.0](#v1130)
  - [Changelog since v1.12.0](#changelog-since-v1120)
  - [Known Issues](#known-issues-2)
  - [Changes by Kind](#changes-by-kind-2)
    - [Deprecation](#deprecation-2)
    - [Features](#features-2)
    - [Bugs](#bugs-2)

# v1.13.2

## Changelog since v1.13.1

## Known Issues

- CSM Operator v1.13.1 may cause pod crashes when PowerFlex response times exceed one second due to premature communication timeout. To resolve this issue, upgrade CSM to v1.14.1 or later through CSM Operator.

## Changes by Kind

### Deprecation

- Support for CSM via Slack will be deprecated on 5/31/2025, aligning with the CSM 1.14 release. All existing Slack channels will be archived from that date onward. Please create a [GitHub Issue](https://github.com/dell/csm/issues) for community support or reach out via the [Dell Support Portal](https://dell.com/support) if you have a valid support contract. For more details, please see our [Dell Support Page](https://dell.github.io/csm-docs/docs/support/).

### Features

### Bugs
- Make probe-retry in Powerflex driver configurable ([#1956](https://github.com/dell/csm/issues/1804))

# v1.13.1

## Changelog since v1.13.0

## Known Issues

## Changes by Kind

### Deprecation

- Support for CSM via Slack will be deprecated on 5/31/2025, aligning with the CSM 1.14 release. All existing Slack channels will be archived from that date onward. Please create a [GitHub Issue](https://github.com/dell/csm/issues) for community support or reach out via the [Dell Support Portal](https://dell.com/support) if you have a valid support contract. For more details, please see our [Dell Support Page](https://dell.github.io/csm-docs/docs/support/).

### Features

### Bugs
- Pods Stuck in Terminating State After PowerFlex CSI Node Pod Restart When Deployments Share Same Node ([#1782](https://github.com/dell/csm/issues/1782))
- Support multiple sidecar headers in the Authorization Proxy Server. ([#1804](https://github.com/dell/csm/issues/1804))

# v1.13.0

## Changelog since v1.12.0

## Known Issues
When using Helm charts to install the PowerMax driver with multiple arrays, the powermax-array-config ConfigMap is incorrectly created, resulting in multiple X_CSI_POWERMAX_ENDPOINT entries. This causes the driver pods to crash with the error "mapping key "X_CSI_POWERMAX_ENDPOINT" already defined". To resolve this issue, you will need to manually edit the ConfigMap powermax-array-config to remove all instances of X_CSI_POWERMAX_ENDPOINT and restart the driver pods (https://github.com/dell/csm/issues/1760).
When restarting a PowerFlex deployment via a command such as 'oc rollout restart', Powerflex CSI node pods will restart. Any deployment whose pods are scheduled on the same node as the restarted pod will be stuck in the Terminating state indefinitely. To resolve this issue, upgrade to v1.13.1 (https://github.com/dell/csm/issues/1782).

## Changes by Kind 

### Deprecation

- Support for CSM via Slack will be deprecated on 5/31/2025, aligning with the CSM 1.14 release. All existing Slack channels will be archived from that date onward. Please create a [GitHub Issue](https://github.com/dell/csm/issues) for community support or reach out via the [Dell Support Portal](https://dell.com/support) if you have a valid support contract. For more details, please see our [Dell Support Page](https://dell.github.io/csm-docs/docs/support/).

### Features

- Deprecation note for slack external support. ([#1679](https://github.com/dell/csm/issues/1679))
- CSI Powerflex must have the ability to connect a subset of the worker nodes to a storage array for multi-array suppport. ([#1613](https://github.com/dell/csm/issues/1613))
- Multi-Availability Zone (AZ) support with multiple storage systems - dedicated storage systems in each AZ. ([#1612](https://github.com/dell/csm/issues/1612))
- Added support for PowerScale 9.10. ([#1611](https://github.com/dell/csm/issues/1611))
- Added Support for PowerStore 4.1. ([#1610](https://github.com/dell/csm/issues/1610))
- Support Kubevirt for CSM modules. ([#1563](https://github.com/dell/csm/issues/1563))
- Added support for Kubernetes 1.32. ([#1561](https://github.com/dell/csm/issues/1561))
- Supporting Openshift 4.18 for CSM.. ([#1560](https://github.com/dell/csm/issues/1560))
- Release CSM 1.13 changes. ([#1559](https://github.com/dell/csm/issues/1559))

### Bugs

- Minimal CR for Powerflex is failing in Csm-operator. ([#1671](https://github.com/dell/csm/issues/1671))
- CSM PowerMax wrong error message. ([#1634](https://github.com/dell/csm/issues/1634))
- CSM deployment minimal file - pulling from quay after updating the image registry. ([#1633](https://github.com/dell/csm/issues/1633))
- csm-metrics-powerstore doesn't start when the PowerStore endpoint is using a DNS name. ([#1632](https://github.com/dell/csm/issues/1632))
- cert-csi CapacityTracking test fails when more than 1 CSI driver is deployed. ([#1504](https://github.com/dell/csm/issues/1504))
- CSM-Authorization unit test race condition with goscaleio v1.18.0. ([#1677](https://github.com/dell/csm/issues/1677))
- PowerMax Driver unable to Delete Replication Group During Replication Operations. ([#1669](https://github.com/dell/csm/issues/1669))
- CSM-Operator is reconciling non CSM pods. ([#1668](https://github.com/dell/csm/issues/1668))
- Labels versions and maintainer update for CSM images. ([#1667](https://github.com/dell/csm/issues/1667))
- Cert-csi is hitting a failure in MultiAttachSuite (MAS). ([#1665](https://github.com/dell/csm/issues/1665))
- Cert-csi modifies ephemeral.volumeAttributes fields causing test failure. ([#1664](https://github.com/dell/csm/issues/1664))
- Pod filesystem not resized while volume gets succesfully expanded. ([#1663](https://github.com/dell/csm/issues/1663))
- A revoked tenant is too tightly coupled to validity of tenant token. ([#1662](https://github.com/dell/csm/issues/1662))
- Latest OPA version (1.0.0) fails to parse Authorization policies. ([#1661](https://github.com/dell/csm/issues/1661))
- Helm installation still check snapshot CRD even though snapshot enabled is set to false. ([#1654](https://github.com/dell/csm/issues/1654))
- Cert-csi is failing for k8s environment. ([#1652](https://github.com/dell/csm/issues/1652))
- PowerMax - X_CSI_IG_MODIFY_HOSTNAME fails to rename a host with same name in different case. ([#1650](https://github.com/dell/csm/issues/1650))
- CSM-Operator: E2E Tests are running with 1 replica count. ([#1648](https://github.com/dell/csm/issues/1648))
- Cannot create PowerMax clones. ([#1644](https://github.com/dell/csm/issues/1644))
- E2E and cert-csi tets are failing. ([#1642](https://github.com/dell/csm/issues/1642))
- NodeGetVolumeStats will cause panic when called w/ an Ephemeral volume. ([#1641](https://github.com/dell/csm/issues/1641))
- CSM PowerFlex entering boot loop when array has long response times. ([#1639](https://github.com/dell/csm/issues/1639))
- CSM Docs Multiple fixes for CSI-Powermax installation. ([#1638](https://github.com/dell/csm/issues/1638))
- Broken links in helm charts readme. ([#1635](https://github.com/dell/csm/issues/1635))
- Documentation requires enhancements for PowerStore multipath.conf. ([#1627](https://github.com/dell/csm/issues/1627))
- PowerScale - handle panic error in ParseNormalizedSnapshotID. ([#1620](https://github.com/dell/csm/issues/1620))
- Prefix with `99-` the CSM-PowerMax MachineConfig sample for multipathing. ([#1618](https://github.com/dell/csm/issues/1618))
- UI is broken for Operator documentation page in the instructions. ([#1615](https://github.com/dell/csm/issues/1615))
- Volume Size Rounding Issue in PowerFlex: Rounds Down Instead of Up for Multiples of 8GB. ([#1608](https://github.com/dell/csm/issues/1608))
- Not able to create CSM using the minimal file, if the Operator deployed from the Operator Hub. ([#1605](https://github.com/dell/csm/issues/1605))
- CSM Operator not deleting the deployment and daemon sets after deleting the CSM. ([#1604](https://github.com/dell/csm/issues/1604))
- CSM Operator Crashing. ([#1603](https://github.com/dell/csm/issues/1603))
- "make install" command is failing for csm-operator. ([#1601](https://github.com/dell/csm/issues/1601))
- Operator e2e scenario for powerscale driver with second set of alternate values is failing in OpenShift cluster. ([#1600](https://github.com/dell/csm/issues/1600))
- Documentation for iSCSI and FC multipathing for PowerStore. ([#1595](https://github.com/dell/csm/issues/1595))
- Remove extra fields from the driver specs when using minimal sample. ([#1594](https://github.com/dell/csm/issues/1594))
- Update the cert-manager version in Powermax Prerequisite. ([#1593](https://github.com/dell/csm/issues/1593))
- Operator e2e scenario for powerflex driver with second set of alternate values is failing in OpenShift cluster. ([#1591](https://github.com/dell/csm/issues/1591))
- Automation for reverseproxy tls secret and  powermax-array-config does not present in E2E. ([#1589](https://github.com/dell/csm/issues/1589))
- Steps for Upgrading Drivers with Dell CSM Operator incorrect or confusing. ([#1588](https://github.com/dell/csm/issues/1588))
- Observability for PowerFlex Creates Too Many Sessions. ([#1587](https://github.com/dell/csm/issues/1587))
- Snapshot from metro volume restore as non-metro even if metro storage class is chosen. ([#1586](https://github.com/dell/csm/issues/1586))
- Stale entries in CSI PowerMax Samples of CSM operator. ([#1585](https://github.com/dell/csm/issues/1585))
- Driver should not be expecting a secret which is not used at all for PowerMax when authorization is enabled. ([#1584](https://github.com/dell/csm/issues/1584))
- CSI-PowerStore Fails to Apply 'mountOptions' Passed in StorageClass. ([#1582](https://github.com/dell/csm/issues/1582))
- Offline bundle doesn't include Authorization Server images. ([#1581](https://github.com/dell/csm/issues/1581))
- Operator offline bundle doesn't prepare registries correctly. ([#1574](https://github.com/dell/csm/issues/1574))
- Apex Navigator for Kubernetes reference be removed from the documentation. ([#1572](https://github.com/dell/csm/issues/1572))
- SubjectAltName needs to be updated in the tls.crt. ([#1571](https://github.com/dell/csm/issues/1571))
- Stale entries in CSM operator samples and helm-charts for PowerMax. ([#1570](https://github.com/dell/csm/issues/1570))
- Unused variable "X_CSI_POWERMAX_ENDPOINT" resulting in driver not to start in PowerMax. ([#1569](https://github.com/dell/csm/issues/1569))
- Examples provided in the secrets of install driver for the Primary Unisphere and Back up Unisphere is lacking clarity in ConfigMap. ([#1568](https://github.com/dell/csm/issues/1568))
- Mode is mentioned incorrectly in the configMap of PowerMax even when it is deployed as a sidecar. ([#1567](https://github.com/dell/csm/issues/1567))
- Inconsistent naming convention of secret is misleading in Installation of PowerMax. ([#1566](https://github.com/dell/csm/issues/1566))
- Documentation for PowerFlex nasName states it is not a required field. ([#1562](https://github.com/dell/csm/issues/1562))
- Prerequisite of multiple secrets which is not necessary with the same cert and key. ([#1557](https://github.com/dell/csm/issues/1557))
- Issue with CSM replication and unable to choose the target cluster certificate. ([#1535](https://github.com/dell/csm/issues/1535))
- snapshot restore failed with Message = failed to get acl entries: Too many links. ([#1514](https://github.com/dell/csm/issues/1514))
- cert-csi documentation clarity. ([#1503](https://github.com/dell/csm/issues/1503))
- Documentation should remove references to CentOS. ([#1467](https://github.com/dell/csm/issues/1467))
- The NVMeCommand constant needs to use full path. ([#1549](https://github.com/dell/csm/issues/1549))
