- [v1.3.1](#v131)
  - [Changelog since v1.3.0](#changelog-since-v130)
  - [Known Issues](#known-issues)
    - [PowerScale Replication: Incorrect quota set on the target PV/directory when Quota enabled](#powerscale-replication-incorrect-quota-set-on-the-target-pvdirectory-when-quota-is-enabled)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)

- [v1.3.0](#v130)
  - [Changelog since v1.2.1](#changelog-since-v121)
  - [Known Issues](#known-issues-1)
    - [PowerScale Replication: Incorrect quota set on the target PV/directory when Quota enabled](#powerscale-replication-incorrect-quota-set-on-the-target-pvdirectory-when-quota-is-enabled-1)
  - [Changes by Kind](#changes-by-kind-1)
    - [Deprecation](#deprecation-1)
    - [Features](#features-1)
    - [Bugs](#bugs-1)

# v1.3.1

## Changelog since v1.3.0 

## Known Issues

### PowerScale Replication: Incorrect quota set on the target PV/directory when Quota is enabled

QuotaScan is not happening correctly causing the SYNCIQ job to fail which is required to create the target PV successfully. In addition, the quota limit size is not being correctly set on the target directories during replication. If a failover is performed in this state, application workloads will encounter an error writing data to the new source volumes (former targets).

### Changes by Kind 

#### Deprecation

- A deprecation note has been added to the [documentation](https://dell.github.io/csm-docs/docs/deployment/csminstaller/) for the CSM Installer, which will be removed in CSM v1.4.0.

#### Features 

- Concurrency enhancements for replicated volumes. ([#550](https://github.com/dell/csm/issues/550))
- Support Volume Expansion of protected PVC with v2.3 driver. ([#553](https://github.com/dell/csm/issues/553))
- Release activities for CSM 1.3.1. ([#590](https://github.com/dell/csm/issues/590))

#### Bugs 

- Volume Attach failed on Node, WWN mismatch. ([#548](https://github.com/dell/csm/issues/548))

# v1.3.0

## Changelog since v1.2.1 

## Known Issues

### PowerScale Replication: Incorrect quota set on the target PV/directory when Quota is enabled

QuotaScan is not happening correctly causing the SYNCIQ job to fail which is required to create the target PV successfully. In addition, the quota limit size is not being correctly set on the target directories during replication. If a failover is performed in this state, application workloads will encounter an error writing data to the new source volumes (former targets).

### Changes by Kind 

#### Deprecation

- A deprecation note has been added to the [documentation](https://dell.github.io/csm-docs/docs/deployment/csminstaller/) for the CSM Installer, which will be removed in CSM v1.4.0.

#### Features 

- CSM Replication support in CSM Operator. ([#270](https://github.com/dell/csm/issues/270))
- CSM Authorization can deployed with Helm. ([#261](https://github.com/dell/csm/issues/261))
- Update Go version to 1.18 for CSM. ([#257](https://github.com/dell/csm/issues/257))
- PowerFlex fsGroup policy support for persistent and CSI Ephemeral volumes. ([#256](https://github.com/dell/csm/issues/256))
- CSM Operator supports PowerScale 2.3. ([#255](https://github.com/dell/csm/issues/255))
- Support for upgrading the CSM Operator and additional enhancements. ([#254](https://github.com/dell/csm/issues/254))
- CSM 1.3 Release Specific Changes. ([#243](https://github.com/dell/csm/issues/243))
- Provide ability to configure multiple Namespaces in a single SRDF Metro RDF group. ([#303](https://github.com/dell/csm/issues/303))
- Enhanced topology control for PowerMax. ([#293](https://github.com/dell/csm/issues/293))
- PowerStore configurable volume attributes. ([#291](https://github.com/dell/csm/issues/291))
- Kubernetes Metadata Enhancement. ([#289](https://github.com/dell/csm/issues/289))
- Add FSGroupPolicy Support in Dell CSI Drivers for PowerMax. ([#285](https://github.com/dell/csm/issues/285))
- Volume to namespace mapping on array for PowerMax. ([#272](https://github.com/dell/csm/issues/272))
- Support for volume group snapshots for PowerStore. ([#268](https://github.com/dell/csm/issues/268))
- Volume Replication Type Upgrade/Downgrade for PowerMax. ([#266](https://github.com/dell/csm/issues/266))
- Common Components for Upgrade/Downgrade volume replication type on CSI provisioned LUN. ([#265](https://github.com/dell/csm/issues/265))
- CSI Driver for PowerScale to allow a Path limit greater than current 128 characters. ([#263](https://github.com/dell/csm/issues/263))
- Add CSM Resiliency support for PowerScale. ([#262](https://github.com/dell/csm/issues/262))
- NVMeoF - FC for PowerStore. ([#260](https://github.com/dell/csm/issues/260))
- Standalone Helm install for the CSI driver for PowerMax. ([#246](https://github.com/dell/csm/issues/246))
- Removing Beta Snapshot sample files. ([#239](https://github.com/dell/csm/issues/239))
- Add support for latest PowerStore array. ([#176](https://github.com/dell/csm/issues/176))
- Monitor CSI Driver node pods failure in CSM for Resiliency so that pods are not scheduled on that node. ([#145](https://github.com/dell/csm/issues/145))

#### Bugs 

- Default nodeSelector/tolerations are not working on k8 1.24 master nodes. ([#319](https://github.com/dell/csm/issues/319))
- Authorization doesn't provide a CLI to list RoleBindings. ([#314](https://github.com/dell/csm/issues/314))
- Verify script is considering OpenShift 4.10 as 4.1. ([#305](https://github.com/dell/csm/issues/305))
- PowerMax: Host group creation is failing for iSCSI configuration. ([#335](https://github.com/dell/csm/issues/335))
- CSI Migrator: VolumeDelete never executes for remote volumes. ([#332](https://github.com/dell/csm/issues/332))
- DiscoverTargets is failing for PowerMax. ([#330](https://github.com/dell/csm/issues/330))
- Expansion of Replicated volume is not supported in CSI-PowerMax. ([#327](https://github.com/dell/csm/issues/327))
- Driver installation throws error with node topology control enabled for Powermax. ([#326](https://github.com/dell/csm/issues/326))
- Resiliency Integration test fails multiple  drivers installed on the same cluster.. ([#325](https://github.com/dell/csm/issues/325))
- Offline bundle for helm based install is not working. ([#324](https://github.com/dell/csm/issues/324))
- Resolve dependabot alerts. ([#323](https://github.com/dell/csm/issues/323))
- Update the branch from master to main in readme for license. ([#320](https://github.com/dell/csm/issues/320))
- Running re-protect on source side puts RG into UNKNOWN state. ([#308](https://github.com/dell/csm/issues/308))
- Retention policy case sensitive and PowerScale RG deletion. ([#307](https://github.com/dell/csm/issues/307))
- Panic occurs in csi-replicator sidecar. ([#306](https://github.com/dell/csm/issues/306))
- CSI driver for Unity StagingTargetPath exceeds size limit for ISCSI PVC. ([#304](https://github.com/dell/csm/issues/304))
- PowerMax: Node pod failed to start when topology control is disabled. ([#300](https://github.com/dell/csm/issues/300))
- PowerStore is not updated with CSI spec it is supporting. ([#295](https://github.com/dell/csm/issues/295))
- Display the correct link state for Empty RG's. ([#286](https://github.com/dell/csm/issues/286))
- CSI Driver for Unity documentation prerequisites should include the same FC and multipath prerequisite steps as found in the PowerStore and PowerMax documentation. ([#284](https://github.com/dell/csm/issues/284))
- Annotation to set the default SC is not beta anymore. ([#281](https://github.com/dell/csm/issues/281))
- New nodes are not recognised by csi-unity driver and Pod NFS mount fails v2.2.0. ([#275](https://github.com/dell/csm/issues/275))
- CODEOWNERS file is missing for https://github.com/dell/csi-powermax. ([#273](https://github.com/dell/csm/issues/273))
- PowerStore CSI appears to look for iSCSI when FC only is being used. ([#252](https://github.com/dell/csm/issues/252))
- CSI reverse proxy logs always display '\n'. ([#251](https://github.com/dell/csm/issues/251))
- Driver installation fails on K8 versions > 1.23.0 using helm. ([#247](https://github.com/dell/csm/issues/247))
- Ext4 filesystem is consuming extra reserved space from the total available size. ([#245](https://github.com/dell/csm/issues/245))
- Resiliency: Occasional failure unmounting Unity volume for raw block devices via iSCSI. ([#237](https://github.com/dell/csm/issues/237))
- CSM Authorization proxy server install fails due to missing container-selinux ([#313](https://github.com/dell/csm/issues/313))
- Permissions on CSM Authorization karavictl and k3s binaries are incorrect ([#277](https://github.com/dell/csm/issues/277))
- CSM Authorization OPA policies fail if there are hosts or DNS entries for "localhost" that don't resolve to 127.0.0.1 ([#321](https://github.com/dell/csm/issues/321))
