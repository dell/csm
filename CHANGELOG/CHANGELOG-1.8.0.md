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

- SLES15 SP4 support in csi powerscale. ([#967](https://github.com/dell/csm/issues/967))
- PowerScale 9.5.0.4 support. ([#950](https://github.com/dell/csm/issues/950))
- CSM for PowerMax file support. ([#861](https://github.com/dell/csm/issues/861))
- CSM support for Openshift 4.13. ([#724](https://github.com/dell/csm/issues/724))
- CSI Unity XT Driver: Add upgrade support to the CSM Operator. ([#955](https://github.com/dell/csm/issues/955))
- Add support for Offline Install of CSM Operator in non OLM environment. ([#939](https://github.com/dell/csm/issues/939))
- Google Anthos 1.15 support  for PowerMax. ([#937](https://github.com/dell/csm/issues/937))
- Use ubi9 micro as base image. ([#922](https://github.com/dell/csm/issues/922))
- Enhancing Unity XT driver to handle API requests after the sessionIdleTimeOut in STIG mode. ([#891](https://github.com/dell/csm/issues/891))
- Enhance GoPowerScale to support PowerScale Terraform Provider. ([#888](https://github.com/dell/csm/issues/888))
- Configurable Volume Attributes use recommended naming convention <prefix>/<name>. ([#879](https://github.com/dell/csm/issues/879))
- CSI 1.5 spec support: Implement Volume Limits. ([#878](https://github.com/dell/csm/issues/878))
- Make standalone helm chart available from helm repository : https://dell.github.io/dell/helm-charts. ([#877](https://github.com/dell/csm/issues/877))
- CSI 1.5 spec support -StorageCapacityTracking. ([#876](https://github.com/dell/csm/issues/876))
- CSI-PowerFlex 4.0 NFS support. ([#763](https://github.com/dell/csm/issues/763))

### Bugs 

- volume expansion is failing with power path. ([#976](https://github.com/dell/csm/issues/976))
- Creating StorageClass for replication failed with unmarshal error. ([#968](https://github.com/dell/csm/issues/968))
- Missing mountpoint in output of lsblk command. ([#966](https://github.com/dell/csm/issues/966))
- Installation Wizard necessary Resiliency fields are commented out. ([#959](https://github.com/dell/csm/issues/959))
- cert-csi help message uses wrong name of "csi-cert". ([#938](https://github.com/dell/csm/issues/938))
- volume-group-snapshot test observes a panic when using "--namespace" parameter. ([#931](https://github.com/dell/csm/issues/931))
- "--attr" of ephemeral-volume performance test doesn't support properties file. ([#930](https://github.com/dell/csm/issues/930))
- PowerStore Replication - Delete RG request hangs. ([#928](https://github.com/dell/csm/issues/928))
- VolumeHealthMetricSuite test failure. ([#923](https://github.com/dell/csm/issues/923))
- Generating report from multiple databases and test runs failure. ([#921](https://github.com/dell/csm/issues/921))
- Remove references to deprecated io/ioutil package. ([#916](https://github.com/dell/csm/issues/916))
- Update Cert-csi documentation for driver certification. ([#914](https://github.com/dell/csm/issues/914))
- Unable to pull podmon image from local repository for offline install. ([#898](https://github.com/dell/csm/issues/898))
- Update CSM Authorization karavictl CLI flag descriptions. ([#895](https://github.com/dell/csm/issues/895))
- CSI driver does not verify iSCSI initiators on the array correctly. ([#849](https://github.com/dell/csm/issues/849))
- Common section for Volume Snapshot Requirements. ([#811](https://github.com/dell/csm/issues/811))
- Powerscale CSI driver RO PVC-from-snapshot wrong zone. ([#487](https://github.com/dell/csm/issues/487))
