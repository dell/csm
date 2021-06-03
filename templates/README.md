## Create application yaml files using ytt templates

### Usage:
To create controller yaml for csi-vxflexos, run:

`ytt -f controller.yaml -f configs/values-vxflexos.yaml -f common/values.yaml -f common/k8s-1.20-values.yaml -f modules/` 

To create controller yaml for csi-vxflexos with podman sidecar, run:

`ytt -f controller.yaml -f configs/values-vxflexos.yaml -f common/values.yaml -f common/k8s-1.20-values.yaml -f modules/ --data-value-yaml podmon.enabled=true`

To write generated yaml files to a directory, run:

` ytt -f controller.yaml -f node.yaml -f csidriver.yaml -f configs/values-vxflexos.yaml -f common/values.yaml -f common/k8s-1.20-values.yaml -f modules/ --output-files=csi-vxflexos-yamls`


### Supported applications
* csi-vxflexos (with optional podman and vgsnapshotter sidecars)
* csi-isilon


To create yaml files for new application, add application specific values.yaml file at config/ and update the module files for different objects if required, for ex to add new sidecar args, volumes etc.