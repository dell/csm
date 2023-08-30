<!--toc-->
- [v1.8.0](#v180)
  - [Changelog since v1.7.1](#changelog-since-v171)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
 

# v1.8.0 

## Changelog since v1.7.1 

## Known Issues 

## Changes by Kind 

### Deprecation 

### Features 

- K8S 1.28 support in CSM 1.8. ([#947](https://github.com/dell/csm/issues/947))
- Add support for Offline Install of CSM Operator in non OLM environment. ([#939](https://github.com/dell/csm/issues/939))
- Set up golangci-lint for all CSM repositories. ([#926](https://github.com/dell/csm/issues/926))
- Use ubi9 micro as base image. ([#922](https://github.com/dell/csm/issues/922))
- Enhancing Unity XT driver to handle API requests after the sessionIdleTimeOut in STIG mode. ([#891](https://github.com/dell/csm/issues/891))
- Make standalone helm chart available from helm repository : https://dell.github.io/dell/helm-charts. ([#877](https://github.com/dell/csm/issues/877))
- CSI 1.5 spec support -StorageCapacityTracking. ([#876](https://github.com/dell/csm/issues/876))
- CSM for PowerMax file support. ([#861](https://github.com/dell/csm/issues/861))
- CSI-PowerFlex 4.0 NFS support. ([#763](https://github.com/dell/csm/issues/763))
- CSM support for Openshift 4.13. ([#724](https://github.com/dell/csm/issues/724))
- Google Anthos 1.15 support  for PowerMax. ([#937](https://github.com/dell/csm/issues/937))
- [BUG]: Unable to pull podmon image from local repository for offline install. ([#898](https://github.com/dell/csm/issues/898))
- Enhance GoPowerScale to support PowerScale Terraform Provider. ([#888](https://github.com/dell/csm/issues/888))
- Configurable Volume Attributes use recommended naming convention <prefix>/<name>. ([#879](https://github.com/dell/csm/issues/879))
- [FEATURE] CSI 1.5 spec support : Implement Volume Limits. ([#878](https://github.com/dell/csm/issues/878))
- Use ubi-micro as base image. ([#790](https://github.com/dell/csm/issues/790))

### Bugs 

- Remove refs to deprecated io/ioutil. ([#916](https://github.com/dell/csm/issues/916))
- Space is not reflecting right on Unity. ([#902](https://github.com/dell/csm/issues/902))
- Unity XT: Volume Mount Hangs. ([#901](https://github.com/dell/csm/issues/901))
- CSI driver does not verify iSCSI initiators on the array correctly. ([#849](https://github.com/dell/csm/issues/849))
- Powerscale CSI driver RO PVC-from-snapshot wrong zone. ([#487](https://github.com/dell/csm/issues/487))
- Powermax : Static provisioning is failing for NFS volume. ([#951](https://github.com/dell/csm/issues/951))
- CSM Operator manual installation missing volumeSnapshot CRDs as a prerequisite. ([#944](https://github.com/dell/csm/issues/944))
- [BUG][cert-csi] : cert-csi  help message use wrong name of "csi-cert". ([#938](https://github.com/dell/csm/issues/938))
- [BUG][cert-csi] : volume-group-snapshot test is observed panic when using "--namespace" parameter. ([#931](https://github.com/dell/csm/issues/931))
- [BUG][cert-csi] :  "--attr" of ephemeral-volume performance test doesn't support properties file. ([#930](https://github.com/dell/csm/issues/930))
- PowerStore Replication - Delete RG request hangs. ([#928](https://github.com/dell/csm/issues/928))
- VolumeHealthMetricSuite test failure. ([#923](https://github.com/dell/csm/issues/923))
- [BUG][cert-csi] :  Generating report from multiple databases and test runs failure. ([#921](https://github.com/dell/csm/issues/921))
- Update Cert-csi documentation for driver certification. ([#914](https://github.com/dell/csm/issues/914))
- Cert Manager should display tooltip about the pre-requisite.. ([#907](https://github.com/dell/csm/issues/907))
- Unable to pull podmon image from local repository for offline install. ([#898](https://github.com/dell/csm/issues/898))
- Documentation - Authorization. ([#895](https://github.com/dell/csm/issues/895))
- Missing nodeSelector and tolerations entry in sample file. ([#890](https://github.com/dell/csm/issues/890))
- Unit tests failing for CSI-PowerMax. ([#887](https://github.com/dell/csm/issues/887))
- Common section for Volume Snapshot Requirements. ([#811](https://github.com/dell/csm/issues/811))
