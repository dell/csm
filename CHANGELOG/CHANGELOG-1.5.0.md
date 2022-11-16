# v1.5.0 

## Changelog since v1.4.0 

## Changes by Kind 

### Features 

- CSM Operator Supports Install of Authorization Proxy Server. ([#511](https://github.com/dell/csm/issues/511))
- CSI-PowerScale: Add support for Standalone Helm charts. ([#506](https://github.com/dell/csm/issues/506))
- CSI Powermax Reverseproxy Mandate. ([#495](https://github.com/dell/csm/issues/495))
- Standardize csi drivers helm installation. ([#494](https://github.com/dell/csm/issues/494))
- CSM 1.5 release specific changes. ([#491](https://github.com/dell/csm/issues/491))
- Support storage capacity tracking in CSI PowerStore driver. ([#483](https://github.com/dell/csm/issues/483))
- Add support for OE 3.2 in CSI PowerStore. ([#482](https://github.com/dell/csm/issues/482))
- CSM support for Openshift 4.11. ([#480](https://github.com/dell/csm/issues/480))
- Include k3s-selinux package as part of CSM Authorization RPM install. ([#409](https://github.com/dell/csm/issues/409))
- Ignore volumeless pods with Resiliency label. ([#493](https://github.com/dell/csm/issues/493))
- CSM Operator: Support install of Observability. ([#488](https://github.com/dell/csm/issues/488))
- Qualify SELinux enablement. ([#394](https://github.com/dell/csm/issues/394))

### Bugs 

- Resolve dependabot reported alerts. ([#540](https://github.com/dell/csm/issues/540))
- ServiceTag is not set in PV volume attributes. ([#538](https://github.com/dell/csm/issues/538))
- CSI PowerScale documentation misses the privileges to user to use CSM replication. ([#516](https://github.com/dell/csm/issues/516))
- PVC fails to resize with the message spec.capacity[storage]: Invalid value: "0": must be greater than zero. ([#507](https://github.com/dell/csm/issues/507))
- Secrets getting regenerated for CSM drivers installed via dell-csi-operator. ([#485](https://github.com/dell/csm/issues/485))
- Documentation Update for Character Count Limits of 64. ([#451](https://github.com/dell/csm/issues/451))
- Node publish is failing with incorrect WWN. ([#542](https://github.com/dell/csm/issues/542))
- PowerMax builds are failing with gomodules error. ([#537](https://github.com/dell/csm/issues/537))
- PowerScale Replication - RESUME doesn't work. ([#535](https://github.com/dell/csm/issues/535))
- karavictl role update can't find existing role in helm deployment. ([#530](https://github.com/dell/csm/issues/530))
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
- CSI-Powerstore: Update the upgrade procedure. ([#484](https://github.com/dell/csm/issues/484))
- step_error: command not found in karavi-observability-install.sh. ([#479](https://github.com/dell/csm/issues/479))
- "repctl cluster inject --use-sa" doesn't work for Kubernetes 1.24 and above. ([#463](https://github.com/dell/csm/issues/463))
