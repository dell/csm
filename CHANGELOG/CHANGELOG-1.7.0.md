- [v1.7.1](#v171)
  - [Changelog since v1.7.0](#changelog-since-v170)
  - [Changes by Kind](#changes-by-kind)
    - [Features](#features)
    
- [v1.7.0](#v170)
  - [Changelog since v1.6.1](#changelog-since-v161)
  - [Known Issues](#known-issues)
    - [CSI Unity XT driver does not verify iSCSI initiators on the array correctly when iSCSI initiator names are not in lowercase](#csi-unity-xt-driver-does-not-verify-iscsi-initiators-on-the-array-correctly-when-iscsi-initiator-names-are-not-in-lowercase) 
    - [CSI Powerstore driver node pods enter CrashLoopBackOff state and provisioning fails](#csi-powerstore-driver-node-pods-enter-crashloopbackoff-state-and-provisioning-fails)
  - [Changes by Kind](#changes-by-kind)
    - [Deprecation](#deprecation) 
    - [Features](#features)
    - [Bugs](#bugs)

# v1.7.1

## Changelog since v1.7.0

## Changes by Kind

### Features 
- Build new image of dellemc/csi-vxflexos tagged as v2.7.1 with latest version of Go. ([#911](https://github.com/dell/csm/issues/911))

# v1.7.0 

## Changelog since v1.6.1 

## Known Issues

### CSI Unity XT driver does not verify iSCSI initiators on the array correctly when iSCSI initiator names are not in lowercase 

After any node reboot, the CSI Unity XT driver pod on that rebooted node goes into a failed state as driver fails to find the iSCSI initiator on the array. The work around is to rename host iSCSI initiators to lowercase and reboot the respective worker node. The CSI Unity XT driver pod will spin off successfully. Example: Rename "iqn.2000-11.com.DEMOWORKERNODE01:1a234b56cd78" to "iqn.2000-11.com.demoworkernode01:1a234b56cd78" in lowercase. 

### CSI Powerstore driver node pods enter CrashLoopBackOff state and provisioning fails

When driver node pods enter CrashLoopBackOff and PVC remains in pending state with one of the following events:
  1. failed to provision volume with StorageClass `<storage-class-name>`: error generating accessibility requirements: no available topology found
  2. waiting for a volume to be created, either by external provisioner "csi-powerstore.dellemc.com" or manually created by system administrator. 

The workaround is to check whether all array details present in the secret file are valid and remove any invalid entries if present. Redeploy the driver.

## Changes by Kind 

#### Deprecation

- CSM for PowerMax linked Proxy mode for [CSI reverse proxy is no longer actively maintained or supported](https://dell.github.io/csm-docs/docs/csidriver/release/powermax/). It will be deprecated in CSM 1.9. It is highly recommended that you use stand alone mode going forward.

### Features 

- CSI Unity XT: Support Unisphere 5.3.0 array. ([#842](https://github.com/dell/csm/issues/842))
- Storage Capacity Tracking for Powerscale. ([#824](https://github.com/dell/csm/issues/824))
- Replication Support for PowerFlex driver in CSM Operator. ([#821](https://github.com/dell/csm/issues/821))
- Add upgrade support of csi-powerstore driver in CSM-Operator. ([#805](https://github.com/dell/csm/issues/805))
- Unity CSI quals and cert for K3s on Debian OS. ([#798](https://github.com/dell/csm/issues/798))
- Validation of authSecret field for Powerstore driver in csm-operator. ([#784](https://github.com/dell/csm/issues/784))
- Support for 'Terrform Provider for PowerFlex 1.1' in goscaleio.. ([#778](https://github.com/dell/csm/issues/778))
- CSM Operator: Add install support for CSI Unity XT driver. ([#756](https://github.com/dell/csm/issues/756))
- CSM 1.7 release specific changes. ([#743](https://github.com/dell/csm/issues/743))
- Allow user to set Quota limit parameters from the PVC request in CSI PowerScale. ([#742](https://github.com/dell/csm/issues/742))
- Add Function to Delete SDC and Change Performance Profile. ([#850](https://github.com/dell/csm/issues/850))
- Update to the latest UBI/UBI Micro images for CSM. ([#843](https://github.com/dell/csm/issues/843))
- PowerMax Support AWS EKS. ([#825](https://github.com/dell/csm/issues/825))
- Powermax : Automate creation of reverse proxy certs. ([#819](https://github.com/dell/csm/issues/819))
- Documentation Enhancement for Replication. ([#818](https://github.com/dell/csm/issues/818))
- Destroy RCG support in repctl. ([#817](https://github.com/dell/csm/issues/817))
- Add Support For PowerFlex Gateway Installer Functions. ([#814](https://github.com/dell/csm/issues/814))
- CSI Powermax - Volumes Not Deleted on Target Array. ([#801](https://github.com/dell/csm/issues/801))
- CSI-PowerMax - Support to mount Block Read-Only PVC. ([#792](https://github.com/dell/csm/issues/792))
- CSM Authorization encryption for secrets in K3S. ([#774](https://github.com/dell/csm/issues/774))
- gopowermax enhancements for PowerMax Terraform provider requirements. ([#770](https://github.com/dell/csm/issues/770))
- CSM Operator: Adds CSI Powermax support. ([#769](https://github.com/dell/csm/issues/769))
- CSM support for Kubernetes 1.27. ([#761](https://github.com/dell/csm/issues/761))
- CSI PowerMax: Support PowerMax v10.0.1. ([#760](https://github.com/dell/csm/issues/760))
- Deprecation notice to remove linked proxy for PowerMax. ([#757](https://github.com/dell/csm/issues/757))
- PowerFlex Replication: Volumes Not Deleted on Target Array. ([#754](https://github.com/dell/csm/issues/754))
- Deprecation of CSI Operator. ([#751](https://github.com/dell/csm/issues/751))
- Add support for host groups for vSphere environment. ([#746](https://github.com/dell/csm/issues/746))
- Migrate image registry from k8s.gcr.io to  registry.k8s.io. ([#744](https://github.com/dell/csm/issues/744))
- CSI PowerStore - Add support for PowerStore Medusa (v3.5) array. ([#735](https://github.com/dell/csm/issues/735))
- CSI PowerMax QoS parameters for throttling performance and bandwidth. ([#726](https://github.com/dell/csm/issues/726))
- CSM Authorization karavictl requires an admin token. ([#725](https://github.com/dell/csm/issues/725))
- CSM Installation Wizard support for CSI PowerScale and PowerFlex and Unity drivers and modules through Helm. ([#698](https://github.com/dell/csm/issues/698))
- CSM Replication: Volumes Not Deleted on Target Array. ([#665](https://github.com/dell/csm/issues/665))

### Bugs 

- Unsupported configurations in support matrix. ([#863](https://github.com/dell/csm/issues/863))
- CSI-Powerflex offline installation is failing during the driver image pull. ([#868](https://github.com/dell/csm/issues/868))
- Known Issues in GitHub releases should be in Troubleshooting section of docs. ([#855](https://github.com/dell/csm/issues/855))
- Fix csm-operator e2e replication and observability tests. ([#853](https://github.com/dell/csm/issues/853))
- CSM-operator: vSphere host id is missing in node manifest. ([#846](https://github.com/dell/csm/issues/846))
- CSI-PowerMax Unit test are failing.. ([#844](https://github.com/dell/csm/issues/844))
- Volume migration from Replication to non-replication is failing. ([#841](https://github.com/dell/csm/issues/841))
- Health-monitor: NodeGetVolumeStats are throwing error. ([#840](https://github.com/dell/csm/issues/840))
- CSM-Operator : Reverse proxy is having incorrect version in sample files. ([#838](https://github.com/dell/csm/issues/838))
- Links are broken in some parts of the Documentation.. ([#836](https://github.com/dell/csm/issues/836))
- PowerFlex parses comments when constructing MDM key. ([#835](https://github.com/dell/csm/issues/835))
- Dellctl should populate all the sidecars and images for the latest release. ([#834](https://github.com/dell/csm/issues/834))
- Standalone binary of cert-csi reports a dependency error. ([#827](https://github.com/dell/csm/issues/827))
- Authorization should have sample CRD for every supported version in csm-operator. ([#826](https://github.com/dell/csm/issues/826))
- Storage Capacity Tracking not working in CSI-PowerStore when installed using CSM Operator. ([#823](https://github.com/dell/csm/issues/823))
- CSM Operator object goes into failed state when deployments are getting scaled down/up. ([#816](https://github.com/dell/csm/issues/816))
- PowerStore CSM Replication Module installation error. ([#815](https://github.com/dell/csm/issues/815))
- CHAP is set to true in the CSI-PowerStore sample file in CSI Operator. ([#812](https://github.com/dell/csm/issues/812))
- CSM Doc improvements for CSI PowerFlex deployment. ([#810](https://github.com/dell/csm/issues/810))
- Remove busybox from Authorization RPM. ([#809](https://github.com/dell/csm/issues/809))
- CSI PowerMax attribute name is mismatched. ([#808](https://github.com/dell/csm/issues/808))
- csm-doc: Update trouble shooting about VM option to resolve mount issue. ([#802](https://github.com/dell/csm/issues/802))
- Improve CSM Operator Authorization documentation. ([#800](https://github.com/dell/csm/issues/800))
- Vsphere creds for vsphere secrets is expected when vsphere enable is set to false. ([#799](https://github.com/dell/csm/issues/799))
- Volume creation is failing with host limits code. ([#797](https://github.com/dell/csm/issues/797))
- csi-install.sh script for csi-powerstore fails with replication CRD even though replication is disabled. ([#795](https://github.com/dell/csm/issues/795))
- Operator install doc inconsistent in test-isilon vs isilon namespace. ([#793](https://github.com/dell/csm/issues/793))
- CSI Driver name. ([#789](https://github.com/dell/csm/issues/789))
- Replication install using Operator does not work. ([#788](https://github.com/dell/csm/issues/788))
- CSM Authorization doesn't write the status code on error for csi-powerscale. ([#787](https://github.com/dell/csm/issues/787))
- Unable to delete application pod when CSI PowerStore is installed using CSM Operator. ([#785](https://github.com/dell/csm/issues/785))
- PowerScale Replication - Target NFS exports are not deleted even though target directories are deleted. ([#782](https://github.com/dell/csm/issues/782))
- Update user requirements for CSI Driver for PowerStore. ([#777](https://github.com/dell/csm/issues/777))
- CSM Operator module unit tests should meet quality criteria. ([#776](https://github.com/dell/csm/issues/776))
- Authorization RPM installation should use nogpgcheck for k3s-selinux package. ([#772](https://github.com/dell/csm/issues/772))
- CSM Authorization - karavictl generate token should output valid yaml. ([#767](https://github.com/dell/csm/issues/767))
- CSI PODMON is tainting the worker node. ([#765](https://github.com/dell/csm/issues/765))
- Failed to create the dell-replication-init  image. ([#758](https://github.com/dell/csm/issues/758))
- Troubleshooting document is missing iscsi related information for PowerMax. ([#750](https://github.com/dell/csm/issues/750))
- Error handling required when CopySnapshot fails for CSI PowerScale. ([#749](https://github.com/dell/csm/issues/749))
- Getting permission denied error when accessing ROX PVC on PoweFlex in OpenShift. ([#745](https://github.com/dell/csm/issues/745))
- CSI-PowerStore: Unable to run e2e test. ([#741](https://github.com/dell/csm/issues/741))
- CSM for Resiliency openshift test required to pass ssh options in scp command. ([#737](https://github.com/dell/csm/issues/737))
- CSM Resiliency GitHub actions produce sporadic failure. ([#733](https://github.com/dell/csm/issues/733))
- gobrick code owner file is containing errors. ([#568](https://github.com/dell/csm/issues/568))
