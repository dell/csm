- [v1.5.1](#v151)
  - [Changelog since v1.5.0](#changelog-since-v150)
  - [Known Issues](#known-issues)
    - [PowerScale Replication: Incorrect quota set on the target PV/directory when Quota enabled](#powerscale-replication-incorrect-quota-set-on-the-target-pvdirectory-when-quota-is-enabled)
  - [Changes by Kind](#changes-by-kind)
    - [Bug Fixes](#bug-fixes)
- [v1.5.0](#v150)
  - [Changelog since v1.4.0](#changelog-since-v140)
  - [Known Issues](#known-issues-1)
    - [PowerScale Replication: Incorrect quota set on the target PV/directory when Quota enabled](#powerscale-replication-incorrect-quota-set-on-the-target-pvdirectory-when-quota-is-enabled-1)
  - [Changes by Kind](#changes-by-kind-1)
    - [Features](#features)
    - [Bug Fixes](#bug-fixes-1)

# v1.5.1 

## Changelog since v1.5.0 

## Known Issues

### PowerScale Replication: Incorrect quota set on the target PV/directory when Quota is enabled

QuotaScan is not happening correctly causing the SYNCIQ job to fail which is required to create the target PV successfully. In addition, the quota limit size is not being correctly set on the target directories during replication. If a failover is performed in this state, application workloads will encounter an error writing data to the new source volumes (former targets).

## Changes by Kind 

### Features 

- API to print a table of the container images. ([#673](https://github.com/dell/csm/issues/673))
- CSM 1.5.1 release specific changes. ([#582](https://github.com/dell/csm/issues/582))
- Show volumes associated with the tenant from the k8s server. ([#408](https://github.com/dell/csm/issues/408))

### Bugs 

- PowerStore CSI driver v2.2 shows "failed to delete mount path" with "device or resource busy". ([#666](https://github.com/dell/csm/issues/666))
- Documentation - PowerStore Update Driver from v2.4 to v2.5 using Helm docs has wrong file name. ([#643](https://github.com/dell/csm/issues/643))
- CSM Authorization installation fails due to a PATH not looking in /usr/local/bin. ([#580](https://github.com/dell/csm/issues/580))

# v1.5.0 

## Changelog since v1.4.0 

## Known Issues

### PowerScale Replication: Incorrect quota set on the target PV/directory when Quota is enabled

QuotaScan is not happening correctly causing the SYNCIQ job to fail which is required to create the target PV successfully. In addition, the quota limit size is not being correctly set on the target directories during replication. If a failover is performed in this state, application workloads will encounter an error writing data to the new source volumes (former targets).

## Changes by Kind 

### Features 

- CSM 1.5 release specific changes. ([#491](https://github.com/dell/csm/issues/491))
- gopowermax support for Terraform PowerMax provider. ([#441](https://github.com/dell/csm/issues/441))
- Automated deployment of SDCs for PowerFlex on RHEL. ([#554](https://github.com/dell/csm/issues/554))
- Scheduled Backups for Application Mobility. ([#551](https://github.com/dell/csm/issues/551))
- Update to the latest UBI/UBI Minimal images for CSM. ([#545](https://github.com/dell/csm/issues/545))
- SLES SP4 Support for CSI PowerFlex driver in CSM. ([#539](https://github.com/dell/csm/issues/539))
- CSI-PowerScale: Add an option to the CSI driver force the client list to be updated even if there are unresolvable host. ([#534](https://github.com/dell/csm/issues/534))
- CSM PowerFlex : QoS parameters for throttling performance and bandwidth of a CSI driver. ([#533](https://github.com/dell/csm/issues/533))
- FC connectivity with virtualized environment. ([#528](https://github.com/dell/csm/issues/528))
- gopowerstore support for Terraform PowerStore provider. ([#520](https://github.com/dell/csm/issues/520))
- CSM Operator Supports Install of Authorization Proxy Server. ([#511](https://github.com/dell/csm/issues/511))
- PowerFlex Read Only Block support. ([#509](https://github.com/dell/csm/issues/509))
- CSI-PowerScale: Add support for Standalone Helm charts. ([#506](https://github.com/dell/csm/issues/506))
- CSI Powermax Reverseproxy Mandate. ([#495](https://github.com/dell/csm/issues/495))
- Standardize csi drivers helm installation. ([#494](https://github.com/dell/csm/issues/494))
- Ignore volumeless pods with Resiliency label. ([#493](https://github.com/dell/csm/issues/493))
- CSM Operator: Support install of Observability. ([#488](https://github.com/dell/csm/issues/488))
- Support storage capacity tracking in CSI PowerStore driver. ([#483](https://github.com/dell/csm/issues/483))
- Add support for OE 3.2 in CSI PowerStore. ([#482](https://github.com/dell/csm/issues/482))
- CSM support for Openshift 4.11. ([#480](https://github.com/dell/csm/issues/480))
- CSM support for Kubernetes 1.25. ([#478](https://github.com/dell/csm/issues/478))
- CSM Operator support to include CSI PowerFlex driver. ([#477](https://github.com/dell/csm/issues/477))
- CSM support for PowerFlex 4.0. ([#476](https://github.com/dell/csm/issues/476))
- Include k3s-selinux package as part of CSM Authorization RPM install. ([#409](https://github.com/dell/csm/issues/409))
- Qualify SELinux enablement. ([#394](https://github.com/dell/csm/issues/394))

### Bugs 

- Remove csireverseproxy check in template. ([#570](https://github.com/dell/csm/issues/570))
- Unmount is failing during node unpublish/unstage calls and volumes are not removed. ([#562](https://github.com/dell/csm/issues/562))
- Deleting metro volume is failing. ([#561](https://github.com/dell/csm/issues/561))
- Verify script is missing latest Openshift version. ([#560](https://github.com/dell/csm/issues/560))
- Typo in description of controllerCount in values.yaml file. ([#559](https://github.com/dell/csm/issues/559))
- Verify script is missing latest Openshift version. ([#557](https://github.com/dell/csm/issues/557))
- Operator : vsphere support for missing on k8 1.21. ([#556](https://github.com/dell/csm/issues/556))
- Replication : Could not add volume in protected SG. ([#552](https://github.com/dell/csm/issues/552))
- Skip verify option used during CSI-Powermax driver installation is asking input from the user for verification of Unisphere REST endpoint support. ([#547](https://github.com/dell/csm/issues/547))
- vSphere : delete volume is not successful. ([#544](https://github.com/dell/csm/issues/544))
- Node publish is failing with incorrect WWN. ([#542](https://github.com/dell/csm/issues/542))
- Resolve dependabot reported alerts. ([#540](https://github.com/dell/csm/issues/540))
- ServiceTag is not set in PV volume attributes. ([#538](https://github.com/dell/csm/issues/538))
- PowerMax builds are failing with gomodules error. ([#537](https://github.com/dell/csm/issues/537))
- PowerScale Replication - RESUME doesn't work. ([#535](https://github.com/dell/csm/issues/535))
- CSM Authorization karavictl role update can't find existing role in helm deployment. ([#530](https://github.com/dell/csm/issues/530))
- CSI-Powerflex doesn't have correct rbac in the Operator. ([#529](https://github.com/dell/csm/issues/529))
- Dynamic values for Director ID and Port ID in gopowermax is not available.. ([#525](https://github.com/dell/csm/issues/525))
- PowerScale Replication - Unplanned failover doesn't work. ([#522](https://github.com/dell/csm/issues/522))
- Observability - Improve Grafana dashboard. ([#519](https://github.com/dell/csm/issues/519))
- PowerScale Replication - Source directory disappears after failover. ([#518](https://github.com/dell/csm/issues/518))
- Known issues are missing in Release Notes documentation for Replication 1.3.0. ([#517](https://github.com/dell/csm/issues/517))
- PowerScale Replication - mount.nfs: Stale file handle after an unplanned failover. ([#515](https://github.com/dell/csm/issues/515))
- PowerScale Replication - Replicated PV has the wrong AzServiceIP. ([#514](https://github.com/dell/csm/issues/514))
- RPM fails to install policies. ([#512](https://github.com/dell/csm/issues/512))
- PowerScale Replication: Differing IsiPaths between Secret and StorageClass cause failure. ([#508](https://github.com/dell/csm/issues/508))
- CSI-PowerFlex RO mount doesn't display correct mount option. ([#503](https://github.com/dell/csm/issues/503))
- ExternalAccess is getting deleted from the PowerStore array when unpublishing the volume. ([#501](https://github.com/dell/csm/issues/501))
- PowerStore CSI driver NVME TCP connectivity, attach volume successful, mount failed with error: csi_sock: connect: connection refused. ([#496](https://github.com/dell/csm/issues/496))
- PowerMax Array support is not updated in CSM docs. ([#489](https://github.com/dell/csm/issues/489))
- Secrets getting regenerated for CSM drivers installed via dell-csi-operator. ([#485](https://github.com/dell/csm/issues/485))
- CSI-Powerstore: Update the upgrade procedure. ([#484](https://github.com/dell/csm/issues/484))
- step_error: command not found in karavi-observability-install.sh. ([#479](https://github.com/dell/csm/issues/479))
- "repctl cluster inject --use-sa" doesn't work for Kubernetes 1.24 and above. ([#463](https://github.com/dell/csm/issues/463))
- Documentation Update for Character Count Limits of 64. ([#451](https://github.com/dell/csm/issues/451))
- Broken link in CSM docs. ([#431](https://github.com/dell/csm/issues/431))
- impossible to install karavi-authorization rpm cause of invalid cross-device link. ([#164](https://github.com/dell/csm/issues/164))
