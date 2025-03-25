<!--toc-->
- [v1.12.1](#v1121)
  - [Changelog since v1.12.0](#changelog-since-v1120)
  - [Changes by Kind](#changes-by-kind)
    - [Bugs](#bugs)
- [v1.12.0](#v1120)
  - [Changelog since v1.11.1](#changelog-since-v1111)
  - [Known Issues](#known-issues)
  - [Changes by Kind](#changes-by-kind-1)
    - [Deprecation](#deprecation)
    - [Features](#features)
    - [Bugs](#bugs-1)

# v1.12.1

## Changelog since v1.12.0

## Changes by Kind

### Bugs

- Support multiple sidecar headers in the Authorization Proxy Server. ([#1804](https://github.com/dell/csm/issues/1804))

# v1.12.0

## Changelog since v1.11.1

## Known Issues
When using CSI-PowerMax with Replication in a multi-cluster setup, the driver on the target cluster fails and logs an error: `“CSI reverseproxy service host or port not found, CSI reverseproxy not installed properly.”`This means the reverseproxy service isn’t set up correctly on the target cluster. You’ll need to manually create the reverseproxy service on the target cluster by following the provided instructions [here](https://dell.github.io/csm-docs/docs/deployment/csmoperator/modules/replication/#configuration-steps)
## Changes by Kind

### Deprecation
- Deprecate Encryption. ([#1517](https://github.com/dell/csm/issues/1517))
- Issue deprecation notice for CSM Volume Group Snapshotter module. ([#1544](https://github.com/dell/csm/issues/1544))
- Deprecate DockerHub - ([#1435](https://github.com/dell/csm/issues/1435))

### Features

- Add Support for KubeVirt. ([#1508](https://github.com/dell/csm/issues/1508))
- Add Support for OpenShift Container Platform (OCP) 4.17. ([#1473](https://github.com/dell/csm/issues/1473))
- Support for Kubernetes 1.31. ([#1472](https://github.com/dell/csm/issues/1472))
- Simplify the CSM Operator deployment. ([#1449](https://github.com/dell/csm/issues/1449))
- CSM 1.12 release specific changes. ([#1435](https://github.com/dell/csm/issues/1435))
- Automatic update of helm charts with latest image version. ([#1414](https://github.com/dell/csm/issues/1414))
- Adding support for PowerMax Magnolia. ([#1410](https://github.com/dell/csm/issues/1410))
- Remove ACC Support. ([#1484](https://github.com/dell/csm/issues/1484))
- Enable Light/Dark Mode Menu in Navbar. ([#1476](https://github.com/dell/csm/issues/1476))
- PowerStore Sync / Metro for Block - CSM Replication. ([#1443](https://github.com/dell/csm/issues/1443))
- Stateless, GitOps, HA enabled deployment of the CSM Authorization proxy server. ([#1281](https://github.com/dell/csm/issues/1281))
- Enable/disable automatic SDC deployment along with driver installation.. ([#663](https://github.com/dell/csm/issues/663))

### Bugs

- Powermax Intergration test failing. ([#1519](https://github.com/dell/csm/issues/1519))
- Dell CSM Installation Issues. ([#1416](https://github.com/dell/csm/issues/1416))
- privTgt mount is lost after vxflexos-node pod restart. ([#1546](https://github.com/dell/csm/issues/1546))
- Helm charts environment variables are missing for powermax-array-config. ([#1543](https://github.com/dell/csm/issues/1543))
- csm-docs support matrix is inconsistent with Unity 5.3.x supported platform. ([#1542](https://github.com/dell/csm/issues/1542))
- CSM Installation Wizard. ([#1540](https://github.com/dell/csm/issues/1540))
- Wrong storage protocol used when multiple PowerStore arrays are defined in secret. ([#1539](https://github.com/dell/csm/issues/1539))
- Host definitions not being created after adding new appliance to secret. ([#1538](https://github.com/dell/csm/issues/1538))
- CSI PowerStore unable to resize NVMe block PVC, even though volume on the array get's resized. ([#1534](https://github.com/dell/csm/issues/1534))
- CSM Operator Will Continually Add Components to Observability. ([#1533](https://github.com/dell/csm/issues/1533))
- CSM-Operator resets dell-replication-controller-config configmap. ([#1531](https://github.com/dell/csm/issues/1531))
- Duplicate host NQNs on nodes with no logs. ([#1530](https://github.com/dell/csm/issues/1530))
- PowerFlex e2e-fsgroup tests are failing. ([#1521](https://github.com/dell/csm/issues/1521))
- iSCSI Linux best practices for PowerStore missing from CSI documentation. ([#1518](https://github.com/dell/csm/issues/1518))
- Missing Node tolerations for resiliency module. ([#1510](https://github.com/dell/csm/issues/1510))
- CSM Operator E2E tests are not passing. ([#1507](https://github.com/dell/csm/issues/1507))
- Fix Gosec error in service.go. ([#1499](https://github.com/dell/csm/issues/1499))
- Incorrect CSI Driver Capability for NFS Volume Cloning in CSM Documentation for PowerFlex. ([#1498](https://github.com/dell/csm/issues/1498))
- Cert-csi Qualification failing for OCP 4.17 rca environment.. ([#1485](https://github.com/dell/csm/issues/1485))
- CSM Operator - Changes to csiDriverSpec does not reflect in CSM state or csidrivers.storage.k8s.io object. ([#1475](https://github.com/dell/csm/issues/1475))
- NVMeTCP Linux requirements and Best Practices need to be documented and/or incorporated into CSI drivers. ([#1465](https://github.com/dell/csm/issues/1465))
- add NVMeTCP connection parameter ctrl-loss-tmo=-1 to implement powerstore best practices. ([#1459](https://github.com/dell/csm/issues/1459))
- CSI-PowerStore Node Prefix is ignored. ([#1458](https://github.com/dell/csm/issues/1458))
- Improve Documentation - Multipath configuration for FC and FC-NVMe attached arrays. ([#1453](https://github.com/dell/csm/issues/1453))
- CSM Installation Wizard page is not rendered properly. ([#1452](https://github.com/dell/csm/issues/1452))
- CSM-operator build fails from disk space issue. ([#1448](https://github.com/dell/csm/issues/1448))
- Gobrick does not clean wwids from /etc/multipath/wwids after removing multipath devices. ([#1447](https://github.com/dell/csm/issues/1447))
- Replication Failover/Reprotect operations has "Error" under State field in the ReplicationGroup. ([#1445](https://github.com/dell/csm/issues/1445))
- Remove mutex locks from interceptors on method calls. ([#1438](https://github.com/dell/csm/issues/1438))
- Access token refresh expiration reverts to 30 seconds. ([#1436](https://github.com/dell/csm/issues/1436))
- Upgrade k8s.io modules in csm-observability module. ([#1431](https://github.com/dell/csm/issues/1431))
- Samples for Cert-CSI documentation is not showing the correct values for storage classes. ([#1428](https://github.com/dell/csm/issues/1428))
- CSM Operator e2e tests: Error in test 3. ([#1427](https://github.com/dell/csm/issues/1427))
- Incorrect Volume Creation Due to Idempotency in CreateVolume. ([#1425](https://github.com/dell/csm/issues/1425))
- Feedback link is not working in CSM-docs. ([#1421](https://github.com/dell/csm/issues/1421))
- Interactive Tutorial unavailable/under maintenance. ([#1419](https://github.com/dell/csm/issues/1419))
- csi-powermax crashed when attempting to unmount volume from node. ([#1418](https://github.com/dell/csm/issues/1418))
- CSM docs home page is not updated with latest matrix. ([#1411](https://github.com/dell/csm/issues/1411))
- csm-authorization helm charts fail against "helm lint". ([#1409](https://github.com/dell/csm/issues/1409))
- CSM Doc's Code block copy button styling is not rendered properly. ([#884](https://github.com/dell/csm/issues/884))
