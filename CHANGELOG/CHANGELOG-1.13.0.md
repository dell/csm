<!--toc-->
- [v1.13.0](#v1130)
  - [Changelog since v1.12.0](#changelog-since-v1120)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
 

# v1.13.0 

## Changelog since v1.12.0 

## Known Issues 

## Changes by Kind 

### Deprecation 

### Features 

- CSI-PowerMax - Mount credentials secret to the reverse-proxy (Customer Ask). ([#1614](https://github.com/dell/csm/issues/1614))
- CSI Powerflex must have the ability to connect a subset of the worker nodes to a storage array for multi-array suppport. ([#1613](https://github.com/dell/csm/issues/1613))
- Multi-Availability Zone (AZ) support with multiple storage systems - dedicated storage systems in each AZ. ([#1612](https://github.com/dell/csm/issues/1612))
- Added support for PowerScale 9.10. ([#1611](https://github.com/dell/csm/issues/1611))
- Added Support for PowerStore 4.1. ([#1610](https://github.com/dell/csm/issues/1610))
- Support Kubevirt for CSM modules. ([#1563](https://github.com/dell/csm/issues/1563))
- Added support for Kubernetes 1.32. ([#1561](https://github.com/dell/csm/issues/1561))
- Supporting Openshift 4.18 for CSM.. ([#1560](https://github.com/dell/csm/issues/1560))
- Release CSM 1.13 changes. ([#1559](https://github.com/dell/csm/issues/1559))

### Bugs 

- UI is broken for Operator documentation page in the instructions. ([#1615](https://github.com/dell/csm/issues/1615))
- Isolation Mechanism for CSI Driver in Shared Storage Pool environments. ([#1606](https://github.com/dell/csm/issues/1606))
- Not able to create CSM using the minimal file, if the Operator deployed from the Operator Hub. ([#1605](https://github.com/dell/csm/issues/1605))
- CSM Operator not deleting the deployment and daemon sets after deleting the CSM. ([#1604](https://github.com/dell/csm/issues/1604))
- CSM Operator Crashing. ([#1603](https://github.com/dell/csm/issues/1603))
- Remove extra fields from the driver specs when using minimal sample. ([#1594](https://github.com/dell/csm/issues/1594))
- Driver should not be expecting a secret which is not used at all for PowerMax when authorization is enabled. ([#1584](https://github.com/dell/csm/issues/1584))
- SubjectAltName needs to be updated in the tls.crt. ([#1571](https://github.com/dell/csm/issues/1571))
- Stale entries in CSM operator samples and helm-charts for PowerMax. ([#1570](https://github.com/dell/csm/issues/1570))
- Inconsistent naming convention of secret is misleading in Installation of PowerMax. ([#1566](https://github.com/dell/csm/issues/1566))
- cert-csi CapacityTracking test fails when more than 1 CSI driver is deployed. ([#1504](https://github.com/dell/csm/issues/1504))
- Broken links in helm charts readme. ([#1635](https://github.com/dell/csm/issues/1635))
- PowerScale - handle panic error in ParseNormalizedSnapshotID. ([#1620](https://github.com/dell/csm/issues/1620))
- Volume Size Rounding Issue in PowerFlex: Rounds Down Instead of Up for Multiples of 8GB. ([#1608](https://github.com/dell/csm/issues/1608))
- "make install" command is failing for csm-operator. ([#1601](https://github.com/dell/csm/issues/1601))
- Operator e2e scenario for powerscale driver with second set of alternate values is failing in OpenShift cluster. ([#1600](https://github.com/dell/csm/issues/1600))
- Documentation for iSCSI and FC multipathing for PowerStore. ([#1595](https://github.com/dell/csm/issues/1595))
- Update the cert-manager version in Powermax Prerequisite. ([#1593](https://github.com/dell/csm/issues/1593))
- Operator e2e scenario for powerflex driver with second set of alternate values is failing in OpenShift cluster. ([#1591](https://github.com/dell/csm/issues/1591))
- Automation for reverseproxy tls secret and  powermax-array-config does not present in E2E. ([#1589](https://github.com/dell/csm/issues/1589))
- Steps for Upgrading Drivers with Dell CSM Operator incorrect or confusing. ([#1588](https://github.com/dell/csm/issues/1588))
- Observability for PowerFlex Creates Too Many Sessions. ([#1587](https://github.com/dell/csm/issues/1587))
- Snapshot from metro volume restore as non-metro even if metro storage class is chosen. ([#1586](https://github.com/dell/csm/issues/1586))
- Stale entries in CSI PowerMax Samples of CSM operator. ([#1585](https://github.com/dell/csm/issues/1585))
- CSI-PowerStore Fails to Apply 'mountOptions' Passed in StorageClass. ([#1582](https://github.com/dell/csm/issues/1582))
- Offline bundle doesn't include Authorization Server images. ([#1581](https://github.com/dell/csm/issues/1581))
- Operator offline bundle doesn't prepare registries correctly. ([#1574](https://github.com/dell/csm/issues/1574))
- Apex Navigator for Kubernetes reference be removed from the documentation. ([#1572](https://github.com/dell/csm/issues/1572))
- Unused variable "X_CSI_POWERMAX_ENDPOINT" resulting in driver not to start in PowerMax. ([#1569](https://github.com/dell/csm/issues/1569))
- Examples provided in the secrets of install driver for the Primary Unisphere and Back up Unisphere is lacking clarity in ConfigMap. ([#1568](https://github.com/dell/csm/issues/1568))
- Mode is mentioned incorrectly in the configMap of PowerMax even when it is deployed as a sidecar. ([#1567](https://github.com/dell/csm/issues/1567))
- Documentation for PowerFlex nasName states it is not a required field. ([#1562](https://github.com/dell/csm/issues/1562))
- Prerequisite of multiple secrets which is not necessary with the same cert and key. ([#1557](https://github.com/dell/csm/issues/1557))
- Issue with CSM replication and unable to choose the target cluster certificate. ([#1535](https://github.com/dell/csm/issues/1535))
- snapshot restore failed with Message = failed to get acl entries: Too many links. ([#1514](https://github.com/dell/csm/issues/1514))
- cert-csi documentation clarity. ([#1503](https://github.com/dell/csm/issues/1503))
- Documentation should remove references to CentOS. ([#1467](https://github.com/dell/csm/issues/1467))
- The NVMeCommand constant needs to use full path. ([#1549](https://github.com/dell/csm/issues/1549))
