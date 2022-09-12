# v1.4.0 

## Changelog since v1.3.0 

## Changes by Kind 

### Features 

- Add installationID in csidriver annotation for PowerFlex CSI. ([#434](https://github.com/dell/csm/issues/434))
- Security improvements for gofsutil. ([#375](https://github.com/dell/csm/issues/375))
- CSI driver for PowerStore handles small requests for filesystems <3GB. ([#365](https://github.com/dell/csm/issues/365))
- Support PowerStore iSCSI volumes with VMware TKG. ([#363](https://github.com/dell/csm/issues/363))
- Documentation enhancements. ([#357](https://github.com/dell/csm/issues/357))
- Standalone helm chart for csi-powerstore. ([#355](https://github.com/dell/csm/issues/355))
- CSM 1.4 Release specific changes. ([#350](https://github.com/dell/csm/issues/350))
- Update Dell CSI Driver storage class sample files to include note regarding defaultFSType value. ([#242](https://github.com/dell/csm/issues/242))
- StorageClass parameter `fsType` is deprecated and should be replaced with `csi.storage.k8s.io/fstype`. ([#188](https://github.com/dell/csm/issues/188))
- Update Container Storage Module images to use the latest UBI/UBI Minimal images. ([#465](https://github.com/dell/csm/issues/465))
- Upgrade the OpenTelemetry Collector for CSM Observability. ([#454](https://github.com/dell/csm/issues/454))
- Set PV/PVC's namespace when using Observability Module. ([#453](https://github.com/dell/csm/issues/453))
- Support PowerScale in CSM Observability. ([#452](https://github.com/dell/csm/issues/452))
- Tech preview for Application Mobility. ([#449](https://github.com/dell/csm/issues/449))
- Encryption for container volumes Tech Preview.. ([#437](https://github.com/dell/csm/issues/437))
- PowerPath support for CSI PowerMax. ([#436](https://github.com/dell/csm/issues/436))
- Juniper support for Powermax. ([#389](https://github.com/dell/csm/issues/389))
- CSM Authorization insecure related entities are renamed to skipCertificateValidation. ([#368](https://github.com/dell/csm/issues/368))
- CSI-Powerscale to add client to only root clients when RO volume created from snapshot and RootClientEnabled. ([#362](https://github.com/dell/csm/issues/362))
- Add support for FsGroupPolicy in Unity XT driver. ([#361](https://github.com/dell/csm/issues/361))

### Bugs 

- Filesystem is not deleted from PowerStore albeit the reclaimPolicy is set to delete when externalAccess is enabled. ([#418](https://github.com/dell/csm/issues/418))
- Authorization: Failing to install k3s in the RPM deployment. ([#461](https://github.com/dell/csm/issues/461))
- Discrepancy  in auto srdf when creating volumes in multiple namespace. ([#440](https://github.com/dell/csm/issues/440))
- The offline installer didn't pull the driver image due to incorrect tag (2.3.0 <> v2.3.0).. ([#435](https://github.com/dell/csm/issues/435))
- Observability Topology: nil pointer error. ([#430](https://github.com/dell/csm/issues/430))
- PowerMax : Failed to find srdf group number for remote volume. ([#420](https://github.com/dell/csm/issues/420))
- PowerScale volumes unable to be created with Helm deployment of CSM Authorization. ([#419](https://github.com/dell/csm/issues/419))
- Authorization CLI documentation does not mention --array-insecure flag when creating or updating storage systems. ([#416](https://github.com/dell/csm/issues/416))
- Authorization: Add documentation for backing up and restoring redis data. ([#410](https://github.com/dell/csm/issues/410))
- CSM Authorization doesn't recognize storage with capital letters. ([#398](https://github.com/dell/csm/issues/398))
- Update Authorization documentation with supported versions of k3s-selinux and container-selinux packages. ([#393](https://github.com/dell/csm/issues/393))
- Using Authorization without dependency on jq. ([#390](https://github.com/dell/csm/issues/390))
- CSI-PowerMax: Expanded size is not reflecting inside the container for File system volumes. ([#386](https://github.com/dell/csm/issues/386))
- Authorization Documentation Improvement. ([#384](https://github.com/dell/csm/issues/384))
- Unit test failing for csm-authorization. ([#382](https://github.com/dell/csm/issues/382))
- PowerMax: Volume expansion is not working for FS volume. ([#378](https://github.com/dell/csm/issues/378))
- Karavictl has incorrect permissions after download. ([#360](https://github.com/dell/csm/issues/360))
- Helm deployment of Authorization denies a valid request path from csi-powerflex. ([#353](https://github.com/dell/csm/issues/353))
