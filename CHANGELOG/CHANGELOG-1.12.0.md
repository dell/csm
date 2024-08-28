<!--toc-->
- [v1.12.0](#v1120)
  - [Changelog since v1.11.0](#changelog-since-v1110)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs)
 

# v1.12.0 

## Changelog since v1.11.0 

## Known Issues 

## Changes by Kind 

### Deprecation 

### Features 

- Unity consistency update to reduce the number of authentication API calls. ([#1415](https://github.com/dell/csm/issues/1415))
- Automatic update of helm charts with latest image version. ([#1414](https://github.com/dell/csm/issues/1414))
- Adding support for PowerMax Magnolia. ([#1410](https://github.com/dell/csm/issues/1410))
- Enable/disable automatic SDC deployment along with driver installation.. ([#663](https://github.com/dell/csm/issues/663))

### Bugs 

- Remove mutex locks from interceptors on method calls. ([#1438](https://github.com/dell/csm/issues/1438))
- Access token refresh expiration reverts to 30 seconds. ([#1436](https://github.com/dell/csm/issues/1436))
- Remove leftover files of old patch driver version in csm-operator. ([#1432](https://github.com/dell/csm/issues/1432))
- Upgrade k8s.io modules in csm-observability module. ([#1431](https://github.com/dell/csm/issues/1431))
- Samples for Cert-CSI documentation should be updated based on the capabilities. ([#1428](https://github.com/dell/csm/issues/1428))
- CSM Operator e2e tests: Error in test 3. ([#1427](https://github.com/dell/csm/issues/1427))
- Add retry step 2 i.e. LinkVolumeToVolume in the subsequent CreateVolume calls. ([#1425](https://github.com/dell/csm/issues/1425))
- Build/unit test execution inconsistencies in README files across repositories. ([#1422](https://github.com/dell/csm/issues/1422))
- csi-powermax crashed when attempting to unmount volume from node. ([#1418](https://github.com/dell/csm/issues/1418))
- Dell CSM Installation Issues. ([#1416](https://github.com/dell/csm/issues/1416))
- Allow to install 2 PowerFlex on a stretched cluster using the CSM Operator or the Helm chart. ([#1413](https://github.com/dell/csm/issues/1413))
- Add samples with symmetrix id for creating volumesnapshotclass referring PowerMax to maintain consistency. ([#1380](https://github.com/dell/csm/issues/1380))
- Feedback link is not working in CSM-docs. ([#1421](https://github.com/dell/csm/issues/1421))
- Interactive Tutorial unavailable/under maintenance. ([#1419](https://github.com/dell/csm/issues/1419))
- CSM docs home page is not updated with latest matrix. ([#1411](https://github.com/dell/csm/issues/1411))
- csm-authorization helm charts fail against "helm lint". ([#1409](https://github.com/dell/csm/issues/1409))
- CSM Doc's Code block copy button styling is not rendered properly. ([#884](https://github.com/dell/csm/issues/884))
