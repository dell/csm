<!--toc-->
- [v1.9.0](#v190)
  - [Changelog since v1.8.0](#changelog-since-v180)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
 

# v1.9.0 

## Changelog since v1.8.0 

## Known Issues 

## Changes by Kind 

### Deprecation 

### Features 

- Update to the latest UBI Micro image for CSM. ([#1031](https://github.com/dell/csm/issues/1031))
- CSM 1.9 release specific changes. ([#1012](https://github.com/dell/csm/issues/1012))
- Dell CSI to Dell CSM Operator Migration Process. ([#996](https://github.com/dell/csm/issues/996))
- Remove linked proxy mode for PowerMax. ([#991](https://github.com/dell/csm/issues/991))
- Add support for CSI Spec 1.6. ([#905](https://github.com/dell/csm/issues/905))

### Bugs 

- Comment out duplicate entries in the sample secret.yaml file. ([#1030](https://github.com/dell/csm/issues/1030))
- Provide more detail about what cert-csi is doing. ([#1027](https://github.com/dell/csm/issues/1027))
- CSM Installation wizard is issuing the warnings that are false positives. ([#1022](https://github.com/dell/csm/issues/1022))
- SDC Rename. ([#1020](https://github.com/dell/csm/issues/1020))
- Run each metrics gatherer in a separate goroutine. ([#1007](https://github.com/dell/csm/issues/1007))
- Too many login sessions in gopowerstore client causes unexpected session termination in UI. ([#1006](https://github.com/dell/csm/issues/1006))
- Allow volume prefix to be set via CSM operator. ([#989](https://github.com/dell/csm/issues/989))
- storageCapacity can be set in unsupported CSI Powermax with CSM Operator. ([#983](https://github.com/dell/csm/issues/983))
- Not able to take volumesnapshots. ([#975](https://github.com/dell/csm/issues/975))
- Missing error check for os.Stat call during volume publish. ([#1014](https://github.com/dell/csm/issues/1014))
- cert-csi invalid path in go.mod prevents installation. ([#1010](https://github.com/dell/csm/issues/1010))
- Cert-CSI from release v1.2.0 downloads wrong version v0.8.1. ([#1009](https://github.com/dell/csm/issues/1009))
- CSM Replication - secret file requirement for both sites not documented. ([#1002](https://github.com/dell/csm/issues/1002))
- Volume health fails because it looks to a wrong path. ([#999](https://github.com/dell/csm/issues/999))
- X_CSI_AUTH_TYPE cannot be set in CSM Operator. ([#990](https://github.com/dell/csm/issues/990))
- CSM Operator fails to install CSM Replication on the remote cluster. ([#988](https://github.com/dell/csm/issues/988))
