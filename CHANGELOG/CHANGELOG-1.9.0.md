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

- CSM 1.9 release specific changes. ([#1012](https://github.com/dell/csm/issues/1012))
- Dell CSI to Dell CSM Operator Migration Process. ([#996](https://github.com/dell/csm/issues/996))
- Add support for CSI Spec 1.6. ([#905](https://github.com/dell/csm/issues/905))

### Bugs 

- Missing error check for os.Stat call during volume publish. ([#1014](https://github.com/dell/csm/issues/1014))
- Run each metrics gatherer in a separate goroutine. ([#1007](https://github.com/dell/csm/issues/1007))
- Too many login sessions in gopowerstore client causes unexpected session termination in UI. ([#1006](https://github.com/dell/csm/issues/1006))
- X_CSI_AUTH_TYPE cannot be set. ([#990](https://github.com/dell/csm/issues/990))
- Allow volume prefix to be set via CSM operator.. ([#989](https://github.com/dell/csm/issues/989))
- storageCapacity can be set in unsupported CSI Powermax with CSM Operator. ([#983](https://github.com/dell/csm/issues/983))
- Not able to take volumesnapshots. ([#975](https://github.com/dell/csm/issues/975))
- cert-csi invalid path in go.mod prevents installation. ([#1010](https://github.com/dell/csm/issues/1010))
- Cert-CSI from release v1.2.0 downloads wrong version v0.8.1. ([#1009](https://github.com/dell/csm/issues/1009))
- Volume health fails because it looks to a wrong path. ([#999](https://github.com/dell/csm/issues/999))
- CSM Operator fails to install CSM Replication on the remote cluster. ([#988](https://github.com/dell/csm/issues/988))
