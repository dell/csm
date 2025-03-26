<!--toc-->
- [v1.14.0](#v1140)
  - [Changelog since v1.12.1](#changelog-since-v1121)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
 

# v1.14.0 

## Changelog since v1.12.1 

## Known Issues 

## Changes by Kind 

### Deprecation 

### Features 

- Improve the support matrix by keeping the same syntax rules. ([#1784](https://github.com/dell/csm/issues/1784))
- CSM PowerFlex - Address API Changes and remove Replication support for PowerFlex 5.0. ([#1759](https://github.com/dell/csm/issues/1759))
- CSM PowerStore - Multiple NAS Servers Support. ([#1758](https://github.com/dell/csm/issues/1758))
- CSM Replication - Test replication failover by creating remote snaps and PVCs/PVs from the snaps. ([#1757](https://github.com/dell/csm/issues/1757))
- CSM Replication - Controller reattach failover PV to PVC automatically on stretched cluster. ([#1756](https://github.com/dell/csm/issues/1756))
- CSM Observability - Grafana Migration from Angular JS to React. ([#1755](https://github.com/dell/csm/issues/1755))
- CSM PowerStore - Host Registration for PowerStore Metro. ([#1753](https://github.com/dell/csm/issues/1753))
- CSM PowerFlex - Expose the SFTP settings to automatically pull the scini.ko kernel module. ([#1752](https://github.com/dell/csm/issues/1752))
- CSM RBAC rules. ([#1751](https://github.com/dell/csm/issues/1751))
- Kubernetes 1.33 Qualification. ([#1750](https://github.com/dell/csm/issues/1750))
- CSM Operator - CSM Operator must manage the CRD only on the K8S cluster where the Operator is deployed. ([#1749](https://github.com/dell/csm/issues/1749))
- CSM PowerMax - Multi-Availability Zone (AZ) support with multiple storage systems - dedicated storage systems in each AZ. ([#1748](https://github.com/dell/csm/issues/1748))
- Release CSM 1.14 changes. ([#1747](https://github.com/dell/csm/issues/1747))
- Add support for  PowerStore 4.2. ([#1746](https://github.com/dell/csm/issues/1746))
- Add support for  PowerFlex 4.8. ([#1745](https://github.com/dell/csm/issues/1745))
- Add support for Unity 5.5. ([#1744](https://github.com/dell/csm/issues/1744))
- Add support for PowerScale 9.11. ([#1743](https://github.com/dell/csm/issues/1743))
- Host Based NFS for File Scalability. ([#1742](https://github.com/dell/csm/issues/1742))
- Deprecate Volume Group Snapshot. ([#1730](https://github.com/dell/csm/issues/1730))
- CSI-PowerMax - Mount credentials secret to the reverse-proxy (Customer Ask). ([#1614](https://github.com/dell/csm/issues/1614))
- Add support for Powermax 10.2. ([#1754](https://github.com/dell/csm/issues/1754))

### Bugs 

- CSI+Rep using Operator for PMAX failing during deployment.. ([#1775](https://github.com/dell/csm/issues/1775))
- gopowerscale CopyIsiVolume* functions ignoring error cases. ([#1773](https://github.com/dell/csm/issues/1773))
- Support multiple sidecar headers in the Authorization Proxy Server. ([#1804](https://github.com/dell/csm/issues/1804))
- Helm Charts release action does not properly support patch releases. ([#1801](https://github.com/dell/csm/issues/1801))
- Clarify backward compatibility statement for Authorization in CSM Documentation. ([#1796](https://github.com/dell/csm/issues/1796))
- Authorization v2 Documentation. ([#1785](https://github.com/dell/csm/issues/1785))
- Pods Stuck in Terminating State After PowerFlex CSI Node Pod Restart When Deployments Share Same Node. ([#1782](https://github.com/dell/csm/issues/1782))
- CSI-PowerFlex - DriverConfigMap is using hard-coded value. ([#1774](https://github.com/dell/csm/issues/1774))
- Broken Links in CSM docs and beta-csm-docs. ([#1772](https://github.com/dell/csm/issues/1772))
- PowerMax node pods are crashing even though the second array is reachable. ([#1769](https://github.com/dell/csm/issues/1769))
- Update the documentation with the OpenShift Virtualization support matrix for CSM modules. ([#1765](https://github.com/dell/csm/issues/1765))
- Cloned PVC remains in a Pending state when cloning a large PVC in PowerScale. ([#1763](https://github.com/dell/csm/issues/1763))
- CSM Operator samples are incomplete. ([#1762](https://github.com/dell/csm/issues/1762))
- Beta Documentation improvements. ([#1761](https://github.com/dell/csm/issues/1761))
- [csi-powermax]: Yaml error in configmap generation. ([#1760](https://github.com/dell/csm/issues/1760))
- CSI PowerFlex does not list volumes from non-default systems. ([#1740](https://github.com/dell/csm/issues/1740))
- failed to provision volume with StorageClass error generating accessibility requirements: no available topology found. ([#1736](https://github.com/dell/csm/issues/1736))
- GitHub Action go-code-tester doesn't report coverage on all packages. ([#1733](https://github.com/dell/csm/issues/1733))
- Updating approveSDC in tenant CR doesn't reflect in backend. ([#1732](https://github.com/dell/csm/issues/1732))
- [CSI-Powerstore] Clarify protocol value use in PV. ([#1729](https://github.com/dell/csm/issues/1729))
- Panic Error During Volume Expansion When Hard Limit is Not Set for CSI PowerScale Driver. ([#1726](https://github.com/dell/csm/issues/1726))
- Scale test fails with powermax nvmetcp protocol when X_CSI_TRANSPORT_PROTOCOL= "". ([#1725](https://github.com/dell/csm/issues/1725))
- cert-csi --longevity does not seem to honor #d option. ([#1720](https://github.com/dell/csm/issues/1720))
- CSM Unity XT protocol/host registration documentation is not clear. ([#1717](https://github.com/dell/csm/issues/1717))
- PowerStore CSI driver version 2.12 - only supports the default interface for iSCSI discovery.. ([#1714](https://github.com/dell/csm/issues/1714))
- Unable to provision PowerMax Metro volumes with replication module not enabled. ([#1711](https://github.com/dell/csm/issues/1711))
- 1.13 documentation | PowerMax | CSI PowerMax Reverse Proxy. ([#1698](https://github.com/dell/csm/issues/1698))
- Auto select protocol makes the node driver to crash. ([#1689](https://github.com/dell/csm/issues/1689))
