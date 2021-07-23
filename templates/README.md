## Create application yaml files using ytt templates

### Usage:
To create controller yaml for any CSI Driver:

Example:

`ytt -f controller.yaml -f node.yaml -f csidriver.yaml -f configs/values-powerflex.yaml -f common/values.yaml -f common/k8s-1.20-values.yaml -f modules/` 

To create controller yaml for any CSI Driver with podmon sidecar:

Example:

`ytt -f controller.yaml -f node.yaml -f csidriver.yaml -f configs/values-powerflex.yaml -f common/values.yaml -f common/k8s-1.20-values.yaml -f configs/values-powerflex-secret.yaml -f modules/ --data-value-yaml podmon.enabled=true`

To write generated yaml files to a directory:

Example:

` ytt -f driversecret.yaml -f controller.yaml -f node.yaml -f csidriver.yaml -f configs/values-powermax.yaml -f common/values.yaml -f common/k8s-1.20-values.yaml -f configs/values-powermax-secret.yaml -f modules/ --output-files=csi-powermax-yamls`


### Supported applications
* CSI-PowerFlex (with optional podman and vgsnapshotter sidecars)
* CSI-PowerScale
* CSI-PowerStore
* CSI-PowerMax
* CSI-Unity

To create yaml files for new application, add application specific values.yaml file at config/ and update the module files for different objects if required, for ex to add new sidecar args, volumes etc.