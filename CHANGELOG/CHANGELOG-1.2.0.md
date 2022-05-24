- [v1.2.1](#v121)
  - [Changelog since v1.2.0](#changelog-since-v120)
  - [Changes by Kind](#changes-by-kind)
    - [Bug Fixes](#bug-fixes)
- [v1.2.0](#v120)
  - [Changelog since v1.1.0](#changelog-since-v110)
  - [Changes by Kind](#changes-by-kind-1)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bug Fixes](#bug-fixes-1)

# v1.2.1

## Changelog since v1.2.0 

### Changes by Kind 

#### Bug Fixes 

- VGS fails to delete all snaps after a group is deleted. ([#292](https://github.com/dell/csm/issues/292))
- PowerStore Grafana dashboard does not populate correctly. ([#279](https://github.com/dell/csm/issues/279))
- Grafana installation script - Prometheus address is incorrect. ([#278](https://github.com/dell/csm/issues/278))
- Observability data values are incorrect. ([#288](https://github.com/dell/csm/issues/288))

# v1.2.0

## Changelog since v1.1.0 

### Changes by Kind 

#### Deprecation

- A deprecation note has been added to the [documentation](https://dell.github.io/csm-docs/docs/deployment/csminstaller/) for the CSM Installer, which will be removed in CSM v1.4.0.

#### Features 

- Support creation of volumes from snapshot with both having different isiPath. ([#227](https://github.com/dell/csm/issues/227))
- CSM Resiliency test enhancements for node affinity. ([#225](https://github.com/dell/csm/issues/225))
- NFSv4 POSIX and ACL support in Dell CSI driver for PowerStore. ([#191](https://github.com/dell/csm/issues/191))
- Update to the latest UBIM image for CSM. ([#183](https://github.com/dell/csm/issues/183))
- Update to the latest Kubernetes CSI Sidecar Container versions for Dell CSI drivers . ([#182](https://github.com/dell/csm/issues/182))
- Add support for new access modes in CSI Spec 1.5 for PowerMax. ([#175](https://github.com/dell/csm/issues/175))
- Add support for Volume Health Monitoring for PowerMax. ([#174](https://github.com/dell/csm/issues/174))
- Add FSGroupPolicy Support in Dell CSI Drivers for PowerScale and PowerStore. ([#167](https://github.com/dell/csm/issues/167))
- CSM Resiliency enhancement to consider pod affinity. ([#165](https://github.com/dell/csm/issues/165))
- Dell CSI Operator Enhancements. ([#161](https://github.com/dell/csm/issues/161))
- New CSM Operator supports PowerScale and Authorization. ([#159](https://github.com/dell/csm/issues/159))
- NVMeoF - TCP for PowerStore. ([#158](https://github.com/dell/csm/issues/158))
- CSM Observability supports upgrade by install script. ([#151](https://github.com/dell/csm/issues/151))
- Update Go version to 1.17 for CSM (across all CSI drivers and modules). ([#149](https://github.com/dell/csm/issues/149))
- Add support for both session-based authentication and basic authentication in CSI-PowerScale. ([#139](https://github.com/dell/csm/issues/139))
- Add support for Kubernetes 1.23. ([#136](https://github.com/dell/csm/issues/136))
- CSM 1.2 Release Specific Changes. ([#128](https://github.com/dell/csm/issues/128))
- CSM Replication Support For PowerScale. ([#116](https://github.com/dell/csm/issues/116))
- Remove CSM for Authorization Sidecar injection in favor of Helm chart deployment. ([#112](https://github.com/dell/csm/issues/112))
- Standalone Helm install for the CSI Unity driver. ([#92](https://github.com/dell/csm/issues/92))
- CSM Resiliency supports evacuation of pods during NoExecute taint on node. ([#87](https://github.com/dell/csm/issues/87))

#### Bug Fixes 

- Dependabot alerts on CSI drivers. ([#178](https://github.com/dell/csm/issues/178))
- Remove trivy scan from build to avoid podman build failures. ([#224](https://github.com/dell/csm/issues/224))
- Removing dirty bit flag from CSI PowerMax 2.2.0 image. ([#222](https://github.com/dell/csm/issues/222))
- Incorrect path for privatemount directory for unity NFS. ([#219](https://github.com/dell/csm/issues/219))
- Documentation in values.yaml related to toleration for running on master node is incorrect - CSI driver for Unity-XT. ([#216](https://github.com/dell/csm/issues/216))
- Volume Health Monitoring section is missing under driver install using Operator. ([#212](https://github.com/dell/csm/issues/212))
- error while creating RO volume from snapshot with different isiPaths. ([#211](https://github.com/dell/csm/issues/211))
- CSI Specification in documentation causes confusion for some users. ([#210](https://github.com/dell/csm/issues/210))
- Leader election timeout flags are not present in operator but present in helm for CSI-Powerscale. ([#209](https://github.com/dell/csm/issues/209))
- Tenant deletion should cancel tenant revokation. ([#208](https://github.com/dell/csm/issues/208))
- CSI request and response ID's are not logged in the driver. ([#206](https://github.com/dell/csm/issues/206))
- CSM authorization / PowerMax / Misleading 401 error on quota violation. ([#205](https://github.com/dell/csm/issues/205))
- Couldn't fence volumes. ([#199](https://github.com/dell/csm/issues/199))
- Scaled down pod and got files from a different volume - Intermittent. ([#198](https://github.com/dell/csm/issues/198))
- Improve deployment documentation for CSM Authorization. ([#197](https://github.com/dell/csm/issues/197))
- Fix the perf issue related to fqdn resolution during export update with clients on isilon. ([#190](https://github.com/dell/csm/issues/190))
- CSI-PowerFlex logs do not contain CSI request IDs on the Request and Response lines. ([#189](https://github.com/dell/csm/issues/189))
- CSM Observability helm charts: Make app.kubernetes.io/name and name consistent. ([#186](https://github.com/dell/csm/issues/186))
- CSI-provisioner container panic issue. ([#180](https://github.com/dell/csm/issues/180))
- Go Mod tidy issues while building the image. ([#172](https://github.com/dell/csm/issues/172))
- Gosec for PowerMax is reporting failure. ([#170](https://github.com/dell/csm/issues/170))
- Verification of secrets repeated twice while installation of driver via helm. ([#168](https://github.com/dell/csm/issues/168))
- CSM Observability documentation is complicated and causing confusion. ([#163](https://github.com/dell/csm/issues/163))
- Documentation has references to using the CSM Installer as recommended method. ([#154](https://github.com/dell/csm/issues/154))
- Unity CSI node driver reports "invalid memory address or nil pointer dereference". ([#152](https://github.com/dell/csm/issues/152))
- Force delete of pod kicks in late (pod in terminating state for a while). ([#148](https://github.com/dell/csm/issues/148))
- CSM Authorization sidecar install fails if k8s worker nodes are not in ~/.ssh/known_hosts. ([#147](https://github.com/dell/csm/issues/147))
- Container is terminated but Pod is stuck in terminating. ([#146](https://github.com/dell/csm/issues/146))
- Dell CSI Operator listed two times after upgrade (1.2.0 + 1.5.0). ([#144](https://github.com/dell/csm/issues/144))
- Failing to create replicated volumes with Integration tests. ([#143](https://github.com/dell/csm/issues/143))
- CSI-PowerScale installation fails when reverse DNS lookup is unavailable. ([#142](https://github.com/dell/csm/issues/142))
- Integration tests for replication is failing with Unsupported replication type. ([#138](https://github.com/dell/csm/issues/138))
- Issues while using PowerFlex secret for Observability. ([#137](https://github.com/dell/csm/issues/137))
- Replication Metro mode is not supported. ([#135](https://github.com/dell/csm/issues/135))
- Documentation improvement recommendations for PowerScale. ([#127](https://github.com/dell/csm/issues/127))
- Driver logs show dirty tag in version. ([#126](https://github.com/dell/csm/issues/126))
- Unit test for csm-metrics-powerstore fails periodically. ([#113](https://github.com/dell/csm/issues/113))

#### Known Issues

- CSM Resiliency: Occasional failure unmounting Unity volume for raw block devices via iSCSI ([#237](https://github.com/dell/csm/issues/237))
