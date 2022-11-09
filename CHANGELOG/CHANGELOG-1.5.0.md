# v1.5.0 

## Changelog since v1.4.0 

## Changes by Kind 

### Features 

- CSM Operator Supports Install of Authorization Proxy Server. ([#511](https://github.com/dell/csm/issues/511))
- CSI-PowerScale: Add support for Standalone Helm charts. ([#506](https://github.com/dell/csm/issues/506))
- CSI Powermax Reverseproxy Mandate. ([#495](https://github.com/dell/csm/issues/495))
- Standardize csi drivers helm installation. ([#494](https://github.com/dell/csm/issues/494))
- CSM 1.5 release specific changes. ([#491](https://github.com/dell/csm/issues/491))
- CSM Operator: Support install of Observability. ([#488](https://github.com/dell/csm/issues/488))
- Support storage capacity tracking in CSI PowerStore driver. ([#483](https://github.com/dell/csm/issues/483))
- Add support for OE 3.2 in CSI PowerStore. ([#482](https://github.com/dell/csm/issues/482))
- CSM support for Openshift 4.11. ([#480](https://github.com/dell/csm/issues/480))
- Include k3s-selinux package as part of CSM Authorization RPM install. ([#409](https://github.com/dell/csm/issues/409))
- Qualify SELinux enablement. ([#394](https://github.com/dell/csm/issues/394))
- File Replication Support For Unity. ([#381](https://github.com/dell/csm/issues/381))
- Ignore volumeless pods with Resiliency label. ([#493](https://github.com/dell/csm/issues/493))

### Bugs 

- karavictl role update can't find existing role in helm deployment. ([#530](https://github.com/dell/csm/issues/530))
- PowerScale Replication - Source directory disappears after failover. ([#518](https://github.com/dell/csm/issues/518))
- CSI PowerScale documentation misses the privileges to user to use CSM replication. ([#516](https://github.com/dell/csm/issues/516))
- PowerScale Replication - mount.nfs: Stale file handle after an unplanned failover. ([#515](https://github.com/dell/csm/issues/515))
- Unity v2.4 MountVolume.WaitForAttach failed. ([#513](https://github.com/dell/csm/issues/513))
- PVC fails to resize with the message spec.capacity[storage]: Invalid value: "0": must be greater than zero. ([#507](https://github.com/dell/csm/issues/507))
- PowerStore CSI driver NVME TCP connectivity, attach volume successful, mount failed with error: csi_sock: connect: connection refused. ([#496](https://github.com/dell/csm/issues/496))
- PowerMax Array support is not updated in CSM docs. ([#489](https://github.com/dell/csm/issues/489))
- Secrets getting regenerated for CSM drivers installed via dell-csi-operator. ([#485](https://github.com/dell/csm/issues/485))
- Documentation Update for Character Count Limits of 64. ([#451](https://github.com/dell/csm/issues/451))
- CSI-Powerflex doesn't have correct rbac in the Operator. ([#529](https://github.com/dell/csm/issues/529))
- Observability - Improve Grafana dashboard. ([#519](https://github.com/dell/csm/issues/519))
- Known issues are missing in Release Notes documentation for Replication 1.3.0. ([#517](https://github.com/dell/csm/issues/517))
- PowerScale Replication - Replicated PV has the wrong AzServiceIP. ([#514](https://github.com/dell/csm/issues/514))
- RPM fails to install policies. ([#512](https://github.com/dell/csm/issues/512))
- PowerScale Replication: Differing IsiPaths between Secret and StorageClass cause failure. ([#508](https://github.com/dell/csm/issues/508))
- CSI-PowerFlex RO mount doesn't display correct mount option. ([#503](https://github.com/dell/csm/issues/503))
- ExternalAccess is getting deleted from the PowerStore array when unpublishing the volume. ([#501](https://github.com/dell/csm/issues/501))
- CSI-Powerstore: Update the upgrade procedure. ([#484](https://github.com/dell/csm/issues/484))
- step_error: command not found in karavi-observability-install.sh. ([#479](https://github.com/dell/csm/issues/479))
