#! /bin/bash

declare -A sidecars_to_repo
sidecars_to_repo[csi-attacher]="https://api.github.com/repos/kubernetes-csi/external-attacher/releases/latest"
sidecars_to_repo[csi-provisioner]="https://api.github.com/repos/kubernetes-csi/external-provisioner/releases/latest"
sidecars_to_repo[csi-snapshotter]="https://api.github.com/repos/kubernetes-csi/external-snapshotter/releases/latest"
sidecars_to_repo[csi-resizer]="https://api.github.com/repos/kubernetes-csi/external-resizer/releases/latest"
sidecars_to_repo[csi-node-driver-registrar]="https://api.github.com/repos/kubernetes-csi/node-driver-registrar/releases/latest"
sidecars_to_repo[csi-external-health-monitor-controller]="https://api.github.com/repos/kubernetes-csi/external-health-monitor/releases/latest"

retrieve_latest_tag() {
    curl -s $1 | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

update_sidecars() {
    sidecar_header="CSI Sidecar versions"
    version_file="config/csm-versions.yaml"

    # Should only be needed on intialization of sidecars
    if ! grep -q "${sidecar_header}" "${version_file}"; then
        echo "" >> $version_file
        echo "# $sidecar_header" >> $version_file
    fi

    for key in "${!sidecars_to_repo[@]}"; do
        echo "Key: $key, Value: ${sidecars_to_repo[$key]}"
        latest_tag=$(retrieve_latest_tag ${sidecars_to_repo[$key]})
        if [ -z "$latest_tag" ]; then
            echo "Failed to retrieve latest image for $key. Try again later."
            exit 1
        fi

        if grep -q $key $version_file; then
            echo "$key found in version, updating version."
            sed -i "s/${key} : .*/${key} : ${latest_tag}/g" $version_file
        else
            echo "$key not found in version, adding."
            echo "$key : $latest_tag" >> $version_file
        fi
    done
}

update_sidecars
