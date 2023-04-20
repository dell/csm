- [v1.6.1](#v161)
  - [Changelog since v1.6.0](#changelog-since-v160)
  - [Changes by Kind](#changes-by-kind)
    - [Bugs](#bugs)
- [v1.6.0](#v160)
  - [Changelog since v1.5.1](#changelog-since-v151)
  - [Known Issues](#known-issues)
    - [PowerScale Replication: Incorrect quota set on the target PV/directory when Quota is enabled](#powerscale-replication-incorrect-quota-set-on-the-target-pvdirectory-when-quota-is-enabled)
  - [Changes by Kind](#changes-by-kind-1)
    - [Features](#features)
    - [Bugs](#bugs-1)

# v1.6.1 

## Changelog since v1.6.0 

## Changes by Kind 

### Bugs 

- PowerScale Replication: Incorrect quota set on the target PV/directory when Quota is enabled. ([#753](https://github.com/dell/csm/issues/753))

# v1.6.0 

## Changelog since v1.5.1 

## Known Issues

### PowerScale Replication: Incorrect quota set on the target PV/directory when Quota is enabled

QuotaScan is not happening correctly causing the SYNCIQ job to fail which is required to create the target PV successfully. In addition, the quota limit size is not being correctly set on the target directories during replication. If a failover is performed in this state, application workloads will encounter an error writing data to the new source volumes (former targets).

## Changes by Kind 

### Features 

- To enable backward compatibility for MKE-3.5.x and K8s 1.21. ([#699](https://github.com/dell/csm/issues/699))
- SLES SP4 Support for CSI Unity XT driver. ([#695](https://github.com/dell/csm/issues/695))
- MKE 3.6.0 support. ([#672](https://github.com/dell/csm/issues/672))
- Ubuntu 22.04 qualification for CSI Powerscale. ([#671](https://github.com/dell/csm/issues/671))
- RKE 1.4.1 support. ([#670](https://github.com/dell/csm/issues/670))
- CSM Operator: Add install support for CSI PowerStore driver. ([#613](https://github.com/dell/csm/issues/613))
- CSM support for Kubernetes 1.26. ([#597](https://github.com/dell/csm/issues/597))
- Add CSM Resiliency support for PowerStore. ([#587](https://github.com/dell/csm/issues/587))
- CSM 1.6 release specific changes. ([#583](https://github.com/dell/csm/issues/583))
- PowerFlex preapproved GUIDs. ([#402](https://github.com/dell/csm/issues/402))
- User friendly name to the PowerFlex. ([#181](https://github.com/dell/csm/issues/181))
- Cleanup powerpath dead paths. ([#669](https://github.com/dell/csm/issues/669))
- Update Go version to 1.20 for CSM 1.6. ([#658](https://github.com/dell/csm/issues/658))
- "access denied by server" while mounting due to NFS refresh on PowerScale. ([#655](https://github.com/dell/csm/issues/655))
- Support restoring of snapshots of metro volumes. ([#652](https://github.com/dell/csm/issues/652))
- Volume cloning of Protected PVCs. ([#646](https://github.com/dell/csm/issues/646))
- Restrict the version of TLS to v1.2 for all requests to CSM authorization proxy server. ([#642](https://github.com/dell/csm/issues/642))
- Cert-csi - Test suite for validating Dell CSI Drivers. ([#628](https://github.com/dell/csm/issues/628))
- CSI PowerFlex - Replication Support. ([#618](https://github.com/dell/csm/issues/618))
- Update to the latest UBI/UBI Minimal images for CSM. ([#612](https://github.com/dell/csm/issues/612))
- CSI Powermax- Snapshot restore to Metro Volumes. ([#609](https://github.com/dell/csm/issues/609))
- CSM Installation Wizard support for CSI PowerStore and PowerMax drivers and modules through Helm. ([#591](https://github.com/dell/csm/issues/591))
- Support APIs for HostGroup in Powermax. ([#588](https://github.com/dell/csm/issues/588))
- Support PowerMax in CSM Observability. ([#586](https://github.com/dell/csm/issues/586))
- Consistent Variable Name for RDM/vSphere. ([#584](https://github.com/dell/csm/issues/584))
- PowerScale replication improvements. ([#573](https://github.com/dell/csm/issues/573))
- PowerScale Replication: Implement Failback functionality. ([#558](https://github.com/dell/csm/issues/558))
- PowerScale Replication - Reprotect should NOT be calling allow_write_revert. ([#532](https://github.com/dell/csm/issues/532))
- Replication APIs to be moved from alpha phase. ([#432](https://github.com/dell/csm/issues/432))

### Bugs 

- gobrick code owner file is containing errors. ([#568](https://github.com/dell/csm/issues/568))
- Multiple bad links in documentation FAQ page and incorrect statements. ([#715](https://github.com/dell/csm/issues/715))
- PowerFlex driver deployment via CSM Operator should work without having the user to create ConfigMap. ([#713](https://github.com/dell/csm/issues/713))
- Migration failed when there are empty MV SG's on source array. ([#705](https://github.com/dell/csm/issues/705))
- Gopowermax unit test Pipeline fails to run tests. ([#701](https://github.com/dell/csm/issues/701))
- Broken link in Container Storage Modules / Deployment / CSM Operator / CSM Modules / Replication. ([#693](https://github.com/dell/csm/issues/693))
- dellctl crashes on a "backup get" when a trailing "/" is added to the namespace. ([#691](https://github.com/dell/csm/issues/691))
- CSM app-mobility can delete restores but they pop back up after 10 seconds.. ([#690](https://github.com/dell/csm/issues/690))
- CSI PowerStore: can't find IP in X_CSI_POWERSTORE_EXTERNAL_ACCESS for NFS provisioning. ([#689](https://github.com/dell/csm/issues/689))
- CSI Powermax fails to create RDF group with free RDF number. ([#688](https://github.com/dell/csm/issues/688))
- vCenter usersname to be updated in secrets.yaml. ([#686](https://github.com/dell/csm/issues/686))
- Operator installation instructions for Replication has broken link. ([#685](https://github.com/dell/csm/issues/685))
- PowerScale Replication: Failback action failing on different environment. ([#677](https://github.com/dell/csm/issues/677))
- Powerstore Multiple iSCSI network support for CSI driver. ([#668](https://github.com/dell/csm/issues/668))
- Deletion of target volume is failing with multiple snapshots. ([#667](https://github.com/dell/csm/issues/667))
- Replication : Creating extra snapshot for idempotent clone operation. ([#657](https://github.com/dell/csm/issues/657))
- CSM Authorization quota of zero should allow infinite use for PowerFlex and PowerMax. ([#654](https://github.com/dell/csm/issues/654))
- Update the optional parameters within angular brackets in upgrade operator page. ([#648](https://github.com/dell/csm/issues/648))
- PowerMax Sample file is wrongly displaying SRDF group auto creation as false. ([#641](https://github.com/dell/csm/issues/641))
- Observability - Improve Grafana dashboards for PowerFlex/PowerStore. ([#640](https://github.com/dell/csm/issues/640))
- CSM Authorization CRD in the CSM Operator doesn't read custom configurations. ([#633](https://github.com/dell/csm/issues/633))
- Powerscale csi driver while connecting smarconnect zone name first letter is getting truncated.. ([#617](https://github.com/dell/csm/issues/617))
- "repctl cluster inject --use-sa" does not work when Replication is installed through `repctl`. ([#600](https://github.com/dell/csm/issues/600))
- PowerStore CSI driver 2.5 create volumes successfully but unable to map volumes to hosts on PowerStore. ([#599](https://github.com/dell/csm/issues/599))
- Delete Snapshot : GetSnapshotInfo() failing with error. ([#593](https://github.com/dell/csm/issues/593))
- vSphere : Can not create VM host object. ([#592](https://github.com/dell/csm/issues/592))
- PowerScale Replication: Artifacts are not properly cleaned after deletion. ([#523](https://github.com/dell/csm/issues/523))
- PVC fails to resize with the message spec.capacity[storage]: Invalid value: "0": must be greater than zero. ([#507](https://github.com/dell/csm/issues/507))
