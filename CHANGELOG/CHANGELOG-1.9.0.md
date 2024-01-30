<!--toc-->
- [v1.9.1](#v191)
  - [Changelog since v1.9.0](#changelog-since-v190)
  - [Changes by Kind](#changes-by-kind)
    - [Bugs](#bugs)
- [v1.9.0](#v190)
  - [Changelog since v1.8.0](#changelog-since-v180)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind-1)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs-1)
 
# v1.9.1

## Changelog since v1.9.0

## Changes by Kind

### Bugs

- Multi Controller defect - sidecars timeout. ([#1110](https://github.com/dell/csm/issues/1110))
- Volumes failing to mount when customer using NVMeTCP on Powerstore. ([#1108](https://github.com/dell/csm/issues/1108))
- Version in Label section of PowerScale v2.9.0 driver is incorrect. ([#1114](https://github.com/dell/csm/issues/1114))
- Operator crashes when deployed from OpenShift with OLM. ([#1117](https://github.com/dell/csm/issues/1117))
- Skip Certificate Validation is not propagated to Authorization module in CSM Operator. ([#1120](https://github.com/dell/csm/issues/1120))
- CSM Operator does not calculate status correctly when module is deployed with driver. ([#1122](https://github.com/dell/csm/issues/1122))

# v1.9.0 

## Changelog since v1.8.0 

## Known Issues 

- For CSM PowerMax, automatic SRDF group creation is failing with "Unable to get Remote Port on SAN for Auto SRDF" on PowerMax 10.1 arrays. As a workaround, create the SRDF Group and add it to the storage class.
- If two separate networks are configured for ISCSI and NVMeTCP, the driver may encounter difficulty identifying the second network (e.g., NVMeTCP). The workaround involves creating a single network on the array to serve both ISCSI and NVMeTCP purposes. ([#1108](https://github.com/dell/csm/issues/1108))
- Standby controller pod is in crashloopbackoff state. Scale down the replica count of the controller pod’s deployment to 1 using ```kubectl scale deployment <deployment_name> --replicas=1 -n <driver_namespace>```. ([#1110](https://github.com/dell/csm/issues/1110))

## Changes by Kind 

### Deprecation 

- The Dell CSI Operator is no longer actively maintained or supported. Dell CSI Operator has been replaced with [Dell CSM Operator](https://dell.github.io/csm-docs/docs/deployment/csmoperator/). If you are currently using Dell CSI Operator, refer to the [operator migration documentation](https://dell.github.io/csm-docs/docs/csidriver/installation/operator/operator_migration/) to migrate from Dell CSI Operator to Dell CSM Operator.
- CSM for PowerMax linked Proxy mode for [CSI reverse proxy is no longer actively maintained or supported](https://dell.github.io/csm-docs/docs/csidriver/release/powermax/). It will be deprecated in CSM 1.9. It is highly recommended that you use stand alone mode going forward.
- The CSM Authorization RPM will be deprecated in a future release. It is highly recommended that you use CSM Authorization [Helm deployment](https://dell.github.io/csm-docs/docs/authorization/deployment/helm/) or [CSM Operator](https://dell.github.io/csm-docs/docs/authorization/deployment/operator/) going forward.

### Features 

- Support For PowerFlex 4.5. ([#1067](https://github.com/dell/csm/issues/1067))
- Support for Openshift 4.14. ([#1066](https://github.com/dell/csm/issues/1066))
- Support for Kubernetes 1.28. ([#947](https://github.com/dell/csm/issues/947))
- CSM PowerMax: Support PowerMax v10.1. ([#1062](https://github.com/dell/csm/issues/1062))
- Update to the latest UBI Micro image for CSM. ([#1031](https://github.com/dell/csm/issues/1031))
- Dell CSI to Dell CSM Operator Migration Process. ([#996](https://github.com/dell/csm/issues/996))
- Remove linked proxy mode for PowerMax. ([#991](https://github.com/dell/csm/issues/991))
- Add support for CSI Spec 1.6. ([#905](https://github.com/dell/csm/issues/905))
- Helm Chart Enhancement - Container Images Configurable in values.yaml. ([#851](https://github.com/dell/csm/issues/851))

### Bugs 

- Documentation links are broken in few places. ([#1072](https://github.com/dell/csm/issues/1072))
- Symmetrix APIs are not getting refreshed. ([#1070](https://github.com/dell/csm/issues/1070))
- CSM Doc page - Update link to PowerStore for Resiliency card. ([#1065](https://github.com/dell/csm/issues/1065))
- Golint is not installing with go get command. ([#1061](https://github.com/dell/csm/issues/1061))
- cert-csi - cannot configure image locations. ([#1059](https://github.com/dell/csm/issues/1059))
- CSI Health monitor for Node missing for CSM PowerFlex in Operator samples. ([#1058](https://github.com/dell/csm/issues/1058))
- CSI Driver - issue with creation volume from 1 of the worker nodes. ([#1057](https://github.com/dell/csm/issues/1057))
- Missing runtime dependencies reference in PowerMax README file.. ([#1056](https://github.com/dell/csm/issues/1056))
- The PowerFlex Dockerfile is incorrectly labeling the version as 2.7.0 for the 2.8.0 version.. ([#1054](https://github.com/dell/csm/issues/1054))
- make gosec is erroring out - Repos PowerMax,PowerStore,PowerScale (gosec is installed). ([#1053](https://github.com/dell/csm/issues/1053))
- make docker command is failing with error. ([#1051](https://github.com/dell/csm/issues/1051))
- NFS Export gets deleted when one pod is deleted from the multiple pods consuming the same PowerFlex RWX NFS volume. ([#1050](https://github.com/dell/csm/issues/1050))
- Is cert-csi expansion expected to successfully run with enableQuota: false on PowerScale?. ([#1046](https://github.com/dell/csm/issues/1046))
- Document instructions update: Either Multi-Path or the Power-Path software should be enabled for PowerMax. ([#1037](https://github.com/dell/csm/issues/1037))
- Comment out duplicate entries in the sample secret.yaml file. ([#1030](https://github.com/dell/csm/issues/1030))
- Provide more detail about what cert-csi is doing. ([#1027](https://github.com/dell/csm/issues/1027))
- CSM Installation wizard is issuing the warnings that are false positives. ([#1022](https://github.com/dell/csm/issues/1022))
- CSI-PowerFlex: SDC Rename fails when configuring multiple arrays in the secret. ([#1020](https://github.com/dell/csm/issues/1020))
- karavi-metrics-powerscale pod gets an segmentation violation error during start. ([#1019](https://github.com/dell/csm/issues/1019))
- Missing error check for os.Stat call during volume publish. ([#1014](https://github.com/dell/csm/issues/1014))
- PowerFlex RWX volume no option to configure the nfs export host access ip address.. ([#1011](https://github.com/dell/csm/issues/1011))
- cert-csi invalid path in go.mod prevents installation. ([#1010](https://github.com/dell/csm/issues/1010))
- Cert-CSI from release v1.2.0 downloads wrong version v0.8.1. ([#1009](https://github.com/dell/csm/issues/1009))
- Too many login sessions in gopowerstore client causes unexpected session termination in UI. ([#1006](https://github.com/dell/csm/issues/1006))
- CSM Replication - secret file requirement for both sites not documented. ([#1002](https://github.com/dell/csm/issues/1002))
- Volume health fails because it looks to a wrong path. ([#999](https://github.com/dell/csm/issues/999))
- X_CSI_AUTH_TYPE cannot be set in CSM Operator. ([#990](https://github.com/dell/csm/issues/990))
- Allow volume prefix to be set via CSM operator. ([#989](https://github.com/dell/csm/issues/989))
- CSM Operator fails to install CSM Replication on the remote cluster. ([#988](https://github.com/dell/csm/issues/988))
- storageCapacity can be set in unsupported CSI Powermax with CSM Operator. ([#983](https://github.com/dell/csm/issues/983))
- Update resources limits for controller-manager to fix OOMKilled error. ([#982](https://github.com/dell/csm/issues/982))
- Not able to take volumesnapshots. ([#975](https://github.com/dell/csm/issues/975))
- Gopowerscale unit test fails. ([#771](https://github.com/dell/csm/issues/771))
