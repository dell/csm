<!--toc-->
- [v1.11.1](#1111)
  - [Changelog since v1.11.0](#changelog-since-v1110)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
- [v1.11.0](#v1110)
  - [Changelog since v1.10.2](#changelog-since-v1102)
  - [Known Issues](#known-issues-1)
  - [Changes by Kind](#changes-by-kind-1)
    - [Deprecation](#deprecation-1)
    - [Features](#features-1)
    - [Bugs](#bugs-1)
 

# v1.11.1

## Changelog since v1.11.0

## Known Issues

## Changes by Kind 

### Deprecation 

### Features 
- Unity consistency update to reduce the number of authentication API calls. ([#1415](https://github.com/dell/csm/issues/1415))

### Bugs
- Mounts using NVMe on PowerStore fails in v2.11 of the driver ([#1469](https://github.com/dell/csm/issues/1469))
- SDC 4.5.2.1 fails to load when deployed on OCP 4.16.x using csm-operator ([#1482](https://github.com/dell/csm/issues/1482))
<br>**Note**: To get the latest 4.5.2.1 SDC image, `ImagePullPolicy` for the `dellemc/sdc:4.5.2.1` image will have to be set to `Always`.

# v1.11.0 

## Changelog since v1.10.2 

## Known Issues

- Controller publish is taking too long to complete. Health monitoring is causing the Unity array to panic by opening multiple sessions. There are error messages in the log `context deadline exceeded`, when health monitoring is enabled. As a workaround, disable volume health monitoring on the node and keep it only at the controller level. Refer [here](https://dell.github.io/csm-docs/docs/csidriver/features/unity/#volume-health-monitoring) for more information about enabling/disabling volume health monitoring.
- Mount of block devices using the NVMe protocols fail. ([#1469](https://github.com/dell/csm/issues/1469)). Use an alternate block protocol (iSCSI, FC) or delay upgrading to the v2.11.0 version of the PowerStore driver and wait for an update of the driver to be published.

## Changes by Kind 

### Deprecation 

### Features 

- Support for Kubernetes 1.30. ([#1400](https://github.com/dell/csm/issues/1400))
- Add Support for OpenShift Container Platform (OCP) 4.16. ([#1359](https://github.com/dell/csm/issues/1359))
- NVMe TCP support for PowerMax. ([#1308](https://github.com/dell/csm/issues/1308))
- Unity 5.4 Support. ([#1399](https://github.com/dell/csm/issues/1399))
- PowerScale  OneFS 9.7 support. ([#1398](https://github.com/dell/csm/issues/1398))
- Observability upgrade is supported in CSM Operator. ([#1397](https://github.com/dell/csm/issues/1397))
- DCM and DN client upgrade is supported in CSM operator. ([#1396](https://github.com/dell/csm/issues/1396))
- Support for PowerFlex 4.6. ([#1358](https://github.com/dell/csm/issues/1358))
- Add Authorization upgrade is supported in CSM Operator. ([#1277](https://github.com/dell/csm/issues/1277))
- CSM Resiliency support for PowerMax. ([#1082](https://github.com/dell/csm/issues/1082))

### Bugs 

- Documentation has broken links to sample files that are no longer available. ([#1392](https://github.com/dell/csm/issues/1392))
- CSI Powermax chooses ISCSI protocol over NVMeTCP. ([#1388](https://github.com/dell/csm/issues/1388))
- Enable static build of repctl. ([#1385](https://github.com/dell/csm/issues/1385))
- Documentation and Release Issues with Cert-CSI Tool. ([#1383](https://github.com/dell/csm/issues/1383))
- Cert-CSI Test Suites "Ephemeral Volumes" is failing. ([#1381](https://github.com/dell/csm/issues/1381))
- Quota capacity limit exceeded. ([#1375](https://github.com/dell/csm/issues/1375))
- Make files in repositories build invalid images. ([#1372](https://github.com/dell/csm/issues/1372))
- API command to check filesystem is taking 20s + causing ControllerUnPublish to take 20+secs. ([#1370](https://github.com/dell/csm/issues/1370))
- Setting large quota in Role causes overflow. ([#1368](https://github.com/dell/csm/issues/1368))
- Support Minimum 3GB Volume Size for NFS in CSI-PowerFlex. ([#1366](https://github.com/dell/csm/issues/1366))
- mkfsFormatOption not working for powerflex. ([#1364](https://github.com/dell/csm/issues/1364))
- Indentation of secret.yaml mentioned on the csm-doc portal for powerflex driver is incorrect.. ([#1355](https://github.com/dell/csm/issues/1355))
- Document update : PowerFlex expecting secret CR as <csm-cr-name>-config in operator. ([#1350](https://github.com/dell/csm/issues/1350))
- karavictl storage create doesn't prompt for storage password. ([#1347](https://github.com/dell/csm/issues/1347))
- Parsing an NVME response fails for list-subsys. ([#1346](https://github.com/dell/csm/issues/1346))
- Helm install of PowerScale does not support snapshots. ([#1340](https://github.com/dell/csm/issues/1340))
- Data loss (DL) when deleting PVC but leaves unusable volumesnapshot and volumesnapshotcontent. ([#1338](https://github.com/dell/csm/issues/1338))
- CSM Replication repctl not supporting static build on OpenSUSE. ([#1330](https://github.com/dell/csm/issues/1330))
- Helm chart is not issuing a warning when installing complex Kubernetes version like Mirantis and alpha/beta versions of Kubernetes. ([#1325](https://github.com/dell/csm/issues/1325))
- PowerScale CSM:  Updating the fsGroupPolicy in the csm is not updating the csidriver. ([#1322](https://github.com/dell/csm/issues/1322))
- Sample file for PowerFlex SDC CR is not updated correctly in the main. ([#1319](https://github.com/dell/csm/issues/1319))
- Link for Dell PowerFlex Deployment Guide is missing in the operator document. ([#1318](https://github.com/dell/csm/issues/1318))
- CSM PowerStore - Remove the RESTAPI code that is not needed. ([#1317](https://github.com/dell/csm/issues/1317))
- PowerScale CSI - Creating PVC from csi snapshot is failing. ([#1316](https://github.com/dell/csm/issues/1316))
- CSI node pod crash after replacing OCP ingress certificate or restarting kubectl service. ([#1310](https://github.com/dell/csm/issues/1310))
- Offline installation documentation appears to be out of date. ([#1307](https://github.com/dell/csm/issues/1307))
- Create volume even if the size is smaller than possible. ([#1305](https://github.com/dell/csm/issues/1305))
- Powerflex snapshots are created as ReadWrite snapshots. ([#1302](https://github.com/dell/csm/issues/1302))
- Missing operator migration page and invalid YAML file version in CSM Docs. ([#1301](https://github.com/dell/csm/issues/1301))
- Images of application mobility velero plugin and controller is not setting the correct image to the latest. ([#1299](https://github.com/dell/csm/issues/1299))
- Fix linter errors in csm-operator. ([#1291](https://github.com/dell/csm/issues/1291))
- CSM docs is having dead links. ([#1289](https://github.com/dell/csm/issues/1289))
- Documentation - RWOP mode has been GAd and it does not need alpha gates anymore. ([#1287](https://github.com/dell/csm/issues/1287))
- unable to install the UNITY driver in NAT Env. ([#1279](https://github.com/dell/csm/issues/1279))
- Installation Wizard creates a 0Byte file when selecting Operator for the installation type. ([#1275](https://github.com/dell/csm/issues/1275))
- Missing entries for Resiliency in installation wizard template. ([#1270](https://github.com/dell/csm/issues/1270))
- Changes in new release of google.golang.org/protobuf is causing compilation issues. ([#1239](https://github.com/dell/csm/issues/1239))
- Missing mountPropagation param for Powermax node template in CSM-Operator. ([#1238](https://github.com/dell/csm/issues/1238))
- Error handling not good in node.go:nodeProbe() and other similar functions. ([#1237](https://github.com/dell/csm/issues/1237))
- Cannot configure export IP for CSI-Unity. ([#1222](https://github.com/dell/csm/issues/1222))
- Issue while Configuring Authorization module with Powermax CSI Driver using Operator. ([#1220](https://github.com/dell/csm/issues/1220))
- Add the helm-charts-version parameter to the install command for all drivers in csm-docs. ([#1218](https://github.com/dell/csm/issues/1218))
- Incorrect Error message in Resiliency Podmon in controllerCleanupPod() func. ([#1216](https://github.com/dell/csm/issues/1216))
- Discrepancy in their secret. ([#1215](https://github.com/dell/csm/issues/1215))
- Doc hyper links in driver Readme is broken. ([#1209](https://github.com/dell/csm/issues/1209))
- Snapshot ingestion procedure for CSI Unity Driver misising. ([#1206](https://github.com/dell/csm/issues/1206))
- Operator doesn't support non-authorization namespace. ([#1205](https://github.com/dell/csm/issues/1205))
- OCP min/max version support. ([#1203](https://github.com/dell/csm/issues/1203))
- CrashLoopBackOff and OOMKilled issue in pod : Dell CSM Operator Manager POD. ([#1200](https://github.com/dell/csm/issues/1200))
- Topology-related node labels are not added automatically. ([#1198](https://github.com/dell/csm/issues/1198))
- Controller Pod keeps restarting due to "Lost connection to CSI driver" error. ([#1188](https://github.com/dell/csm/issues/1188))
