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

- Remove ACC Support. ([#1484](https://github.com/dell/csm/issues/1484))
- Simplify the CSM Operator deployment. ([#1449](https://github.com/dell/csm/issues/1449))
- PowerStore Sync / Metro for Block - CSM Replication. ([#1443](https://github.com/dell/csm/issues/1443))
- Automatic update of helm charts with latest image version. ([#1414](https://github.com/dell/csm/issues/1414))
- Adding support for PowerMax Magnolia. ([#1410](https://github.com/dell/csm/issues/1410))
- Stateless, GitOps, HA enabled deployment of the CSM Authorization proxy server. ([#1281](https://github.com/dell/csm/issues/1281))
- Enable/disable automatic SDC deployment along with driver installation.. ([#663](https://github.com/dell/csm/issues/663))

### Bugs 

- Cert-csi Qualification failing for OCP 4.17 rca environment.. ([#1485](https://github.com/dell/csm/issues/1485))
- Gobrick does not clean wwids from /etc/multipath/wwids after removing multipath devices. ([#1447](https://github.com/dell/csm/issues/1447))
- Upgrade k8s.io modules in csm-observability module. ([#1431](https://github.com/dell/csm/issues/1431))
- CSM Operator e2e tests: Error in test 3. ([#1427](https://github.com/dell/csm/issues/1427))
- csi-powermax crashed when attempting to unmount volume from node. ([#1418](https://github.com/dell/csm/issues/1418))
- Dell CSM Installation Issues. ([#1416](https://github.com/dell/csm/issues/1416))
- Allow to install 2 PowerFlex on a stretched cluster using the CSM Operator or the Helm chart. ([#1413](https://github.com/dell/csm/issues/1413))
- CSM Operator - Changes to csiDriverSpec does not reflect in CSM state or csidrivers.storage.k8s.io object. ([#1475](https://github.com/dell/csm/issues/1475))
- CSI-PowerStore Node Prefix is ignored. ([#1458](https://github.com/dell/csm/issues/1458))
- CSM Installation Wizard page is not rendered properly. ([#1452](https://github.com/dell/csm/issues/1452))
- Replication Failover/Reprotect operations has "Error" under State field in the ReplicationGroup. ([#1445](https://github.com/dell/csm/issues/1445))
- Remove mutex locks from interceptors on method calls. ([#1438](https://github.com/dell/csm/issues/1438))
- Access token refresh expiration reverts to 30 seconds. ([#1436](https://github.com/dell/csm/issues/1436))
- Samples for Cert-CSI documentation is not showing the correct values for storage classes. ([#1428](https://github.com/dell/csm/issues/1428))
- Incorrect Volume Creation Due to Idempotency in CreateVolume. ([#1425](https://github.com/dell/csm/issues/1425))
- Feedback link is not working in CSM-docs. ([#1421](https://github.com/dell/csm/issues/1421))
- Interactive Tutorial unavailable/under maintenance. ([#1419](https://github.com/dell/csm/issues/1419))
- CSM docs home page is not updated with latest matrix. ([#1411](https://github.com/dell/csm/issues/1411))
- csm-authorization helm charts fail against "helm lint". ([#1409](https://github.com/dell/csm/issues/1409))
- CSM Doc's Code block copy button styling is not rendered properly. ([#884](https://github.com/dell/csm/issues/884))
