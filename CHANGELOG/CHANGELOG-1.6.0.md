# v1.6.0 

## Changelog since v1.5.0 

## Changes by Kind 

### Features 

- MKE 3.6.0 support. ([#672](https://github.com/dell/csm/issues/672))
- RKE 1.4.1 support. ([#670](https://github.com/dell/csm/issues/670))
- Update Go version to 1.20 for CSM 1.6. ([#658](https://github.com/dell/csm/issues/658))
- Volume cloning of Protected PVCs. ([#646](https://github.com/dell/csm/issues/646))
- Cert-csi - Test suite for validating Dell CSI Drivers. ([#628](https://github.com/dell/csm/issues/628))
- CSI PowerFlex - Replication Support. ([#618](https://github.com/dell/csm/issues/618))
- CSM Operator: Add install support for CSI PowerStore driver. ([#613](https://github.com/dell/csm/issues/613))
- Update to the latest UBI/UBI Minimal images for CSM. ([#612](https://github.com/dell/csm/issues/612))
- CSM support for Kubernetes 1.26. ([#597](https://github.com/dell/csm/issues/597))
- CSM Installation Wizard support for CSI PowerStore and PowerMax drivers and modules through Helm. ([#591](https://github.com/dell/csm/issues/591))
- Support APIs for HostGroup in Powermax. ([#588](https://github.com/dell/csm/issues/588))
- Add CSM Resiliency support for PowerStore. ([#587](https://github.com/dell/csm/issues/587))
- Consistent Variable Name for RDM/vSphere. ([#584](https://github.com/dell/csm/issues/584))
- CSM 1.6 release specific changes. ([#583](https://github.com/dell/csm/issues/583))
- PowerFlex preapproved GUIDs. ([#402](https://github.com/dell/csm/issues/402))
- User friendly name to the PowerFlex. ([#181](https://github.com/dell/csm/issues/181))
- Cleanup powerpath dead paths. ([#669](https://github.com/dell/csm/issues/669))
- CSI Powermax- Snapshot restore to Metro Volumes. ([#609](https://github.com/dell/csm/issues/609))
- Support PowerMax in CSM Observability. ([#586](https://github.com/dell/csm/issues/586))
- PowerScale replication improvements. ([#573](https://github.com/dell/csm/issues/573))
- PowerScale Replication: Implement Failback functionality. ([#558](https://github.com/dell/csm/issues/558))
- PowerScale Replication - Reprotect should NOT be calling allow_write_revert. ([#532](https://github.com/dell/csm/issues/532))
- Replication APIs to be moved from alpha phase. ([#432](https://github.com/dell/csm/issues/432))

### Bugs 

- PowerScale Replication: Failback action failing on different environment. ([#677](https://github.com/dell/csm/issues/677))
- Replication : Creating extra snapshot for idempotent clone operation. ([#657](https://github.com/dell/csm/issues/657))
- CSM Authorization quota of zero should allow infinite use for PowerFlex and PowerMax. ([#654](https://github.com/dell/csm/issues/654))
- Update the optional parameters within angular brackets in upgrade operator page. ([#648](https://github.com/dell/csm/issues/648))
- gobrick code owner file is containing errors. ([#568](https://github.com/dell/csm/issues/568))
- PVC fails to resize with the message spec.capacity[storage]: Invalid value: "0": must be greater than zero. ([#507](https://github.com/dell/csm/issues/507))
- Powerstore Multiple iSCSI network support for CSI driver. ([#668](https://github.com/dell/csm/issues/668))
- Deletion of target volume is failing with multiple snapshots. ([#667](https://github.com/dell/csm/issues/667))
- PowerMax Sample file is wrongly displaying SRDF group auto creation as false. ([#641](https://github.com/dell/csm/issues/641))
- Observability - Improve Grafana dashboards for PowerFlex/PowerStore. ([#640](https://github.com/dell/csm/issues/640))
- CSM Authorization CRD in the CSM Operator doesn't read custom configurations. ([#633](https://github.com/dell/csm/issues/633))
- Powerscale csi driver while connecting smarconnect zone name first letter is getting truncated.. ([#617](https://github.com/dell/csm/issues/617))
- "repctl cluster inject --use-sa" does not work when Replication is installed through `repctl`. ([#600](https://github.com/dell/csm/issues/600))
- PowerStore CSI driver 2.5 create volumes successfully but unable to map volumes to hosts on PowerStore. ([#599](https://github.com/dell/csm/issues/599))
- Delete Snapshot : GetSnapshotInfo() failing with error. ([#593](https://github.com/dell/csm/issues/593))
- vSphere : Can not create VM host object. ([#592](https://github.com/dell/csm/issues/592))
- PowerScale Replication: Artifacts are not properly cleaned after deletion. ([#523](https://github.com/dell/csm/issues/523))
