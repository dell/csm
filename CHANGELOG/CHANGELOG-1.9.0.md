<!--toc-->
- [v1.9.2](#v192)
  - [Changelog since v1.9.1](#changelog-since-v191)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Bugs](#bugs)
- [v1.9.1](#v191)
  - [Changelog since v1.9.0](#changelog-since-v190)
  - [Known Issues](#known-issues-1)
  - [Changes by Kind](#changes-by-kind-1)
    - [Bugs](#bugs-1)
- [v1.9.0](#v190)
  - [Changelog since v1.8.0](#changelog-since-v180)
  - [Known Issues](#known-issues-2)
  - [Changes by Kind](#changes-by-kind-2)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs-2)
 
# v1.9.2

## Changelog since v1.9.1

## Known Issues

- The status field of a csm object as deployed by CSM Operator may, in limited cases, display a "Failed" status for a successful deployment. As a workaround, the deployment is still usable as long as all pods are running/healthy.
- The status calculation done for the csm object associated with the Authorization Proxy Server when deployed with CSM Operator assumes that the proxy server will be deployed in the "authorization" namespace. If a different namespace is used, the status will stay in the failed state, even though the deployment is healthy. As a workaround, we recommend using the "authorization" namespace for the proxy server. If this is not possible, the health of the deployment can be verified by checking the status of all the pods rather than by checking the status field.

## Changes by Kind

### Bugs

- CSM Operator doesn't apply fSGroupPolicy value to CSIDriver Object. ([#1103](https://github.com/dell/csm/issues/1103))
- CSM Operator does not calculate status correctly when a driver is deployed by itself. ([#1130](https://github.com/dell/csm/issues/1130))
- CSM Operator does not calculate status correctly when application-mobility is deployed by itself. ([#1133](https://github.com/dell/csm/issues/1133))
- CSM Operator intermittently does not calculate status correctly when deploying a driver. ([#1137](https://github.com/dell/csm/issues/1137))
- CSM Operator does not calculate status correctly when deploying the authorization proxy server. ([#1143](https://github.com/dell/csm/issues/1143))
- CSM Operator does not calculate status correctly when deploying observability with csi-powerscale. ([#1146](https://github.com/dell/csm/issues/1146))
- CSM Operator labels csm objects with CSMVersion 1.8.0, an old version. ([#1147](https://github.com/dell/csm/issues/1147))

# v1.9.1

## Changelog since v1.9.0

## Known Issues 

- For CSM Operator released in CSM v1.9.1, a plain driver install (no modules) will always be marked as failed in the CSM status even when it succeeds. As a workaround, the driver deployment is still usable as long as all the pods are running/healthy.
- For CSM Operator released in CSM v1.9.1, a standalone install of application-mobility (not as a module under the driver CSM) will always be marked as failed in the CSM status, even when it succeeds. This is because the operator is looking for the wrong daemonset label to confirm the deployment. As a workaround, the module is still usable as long as all the pods are running/healthy.
- For CSM Operator released in CSM v1.9.1, a driver install will rarely (~2% of the time) have a csm object stuck in a failed state for over an hour even though the deployment succeeds. This is due to a race condition in the status update logic. As a workaround, the driver is still usable as long as all the pods are running/healthy.
- For CSM Operator released in CSM v1.9.1, the authorization proxy server csm object status will always be failed, even when it succeeds. This is because the operator is looking for a daemonset status when the authorization proxy server deployment does not have a daemonset. As a workaround, the module is still usable as long as all the pods are running/healthy.
- For CSM Operator released in CSM v1.9.1, an install of csi-powerscale with observability will always be marked as failed in the csm object status, even when it succeeds. This is because the operator is looking for a legacy name of isilon in the status check. As a workaround, the module is still usable as long as all the pods are running/healthy.
- For csm objects created by the CSM Operator, the CSMVersion label value is v1.8.0 when it should be v1.9.1. As a workaround, the CSM version can be double-checked by checking the operator version -- v1.4.1 operator corresponds to CSM v1.9.1.
- The status field of a csm object as deployed by CSM Operator may, in limited cases, display a "Failed" status for a successful deployment. As a workaround, the deployment is still usable as long as all pods are running/healthy.

## Changes by Kind

### Bugs

- Multi Controller defect - sidecars timeout. ([#1110](https://github.com/dell/csm/issues/1110))
- Volumes failing to mount when customer using NVMeTCP on Powerstore. ([#1108](https://github.com/dell/csm/issues/1108))
- Operator crashes when deployed from OpenShift with OLM. ([#1117](https://github.com/dell/csm/issues/1117))
- Skip Certificate Validation is not propagated to Authorization module in CSM Operator. ([#1120](https://github.com/dell/csm/issues/1120))
- CSM Operator does not calculate status correctly when module is deployed with driver. ([#1122](https://github.com/dell/csm/issues/1122))

# v1.9.0 

## Changelog since v1.8.0 

## Known Issues 

- For CSM PowerMax, automatic SRDF group creation is failing with "Unable to get Remote Port on SAN for Auto SRDF" on PowerMax 10.1 arrays. As a workaround, create the SRDF Group and add it to the storage class.
- For CSM Operator released in CSM v1.9.0, a driver install will rarely (~2% of the time) have a csm object stuck in a failed state for over an hour even though the deployment succeeds. This is due to a race condition in the status update logic.
- For csm objects created by the CSM Operator, the CSMVersion label value is v1.8.0 when it should be v1.9.0. As a workaround, the CSM version can be double-checked by checking the operator version -- v1.4.0 operator corresponds to CSM v1.9.0.
- The status field of a csm object as deployed by CSM Operator may, in limited cases, display a "Failed" status for a successful deployment. As a workaround, the deployment is still usable as long as all pods are running/healthy.
  
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
