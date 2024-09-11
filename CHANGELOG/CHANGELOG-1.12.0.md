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

- CSM Operator - Need for a simple and minimal manifest file using existing Operator CRD. ([#1449](https://github.com/dell/csm/issues/1449))
- PowerStore Sync / Metro for Block - CSM Replication. ([#1443](https://github.com/dell/csm/issues/1443))
- Automatic update of helm charts with latest image version. ([#1414](https://github.com/dell/csm/issues/1414))
- Adding support for PowerMax Magnolia. ([#1410](https://github.com/dell/csm/issues/1410))
- Enable/disable automatic SDC deployment along with driver installation.. ([#663](https://github.com/dell/csm/issues/663))
- Unity consistency update to reduce the number of authentication API calls. ([#1415](https://github.com/dell/csm/issues/1415))

### Bugs 

- CSI-PowerStore Node Prefix is ignored. ([#1458](https://github.com/dell/csm/issues/1458))
- Gobrick does not clean wwids from /etc/multipath/wwids after removing multipath devices. ([#1447](https://github.com/dell/csm/issues/1447))
- Remove leftover files of old patch driver version in csm-operator. ([#1432](https://github.com/dell/csm/issues/1432))
- Upgrade k8s.io modules in csm-observability module. ([#1431](https://github.com/dell/csm/issues/1431))
- CSM Operator e2e tests: Error in test 3. ([#1427](https://github.com/dell/csm/issues/1427))
- csi-powermax crashed when attempting to unmount volume from node. ([#1418](https://github.com/dell/csm/issues/1418))
- Dell CSM Installation Issues. ([#1416](https://github.com/dell/csm/issues/1416))
- Allow to install 2 PowerFlex on a stretched cluster using the CSM Operator or the Helm chart. ([#1413](https://github.com/dell/csm/issues/1413))
- Samples for volumesnapshot class for taking snapshots on target array are not present in repositories other than pmax. ([#1380](https://github.com/dell/csm/issues/1380))
- CSM Installation Wizard page is not rendered properly. ([#1452](https://github.com/dell/csm/issues/1452))
- Remove mutex locks from interceptors on method calls. ([#1438](https://github.com/dell/csm/issues/1438))
- Access token refresh expiration reverts to 30 seconds. ([#1436](https://github.com/dell/csm/issues/1436))
- Samples for Cert-CSI documentation is not showing the correct values for storage classes. ([#1428](https://github.com/dell/csm/issues/1428))
- Incorrect Volume Creation Due to Idempotency in CreateVolume. ([#1425](https://github.com/dell/csm/issues/1425))
- Feedback link is not working in CSM-docs. ([#1421](https://github.com/dell/csm/issues/1421))
- Interactive Tutorial unavailable/under maintenance. ([#1419](https://github.com/dell/csm/issues/1419))
- CSM docs home page is not updated with latest matrix. ([#1411](https://github.com/dell/csm/issues/1411))
- csm-authorization helm charts fail against "helm lint". ([#1409](https://github.com/dell/csm/issues/1409))
- CSM Doc's Code block copy button styling is not rendered properly. ([#884](https://github.com/dell/csm/issues/884))
