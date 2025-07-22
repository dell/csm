#! /bin/bash

declare -A latest_images
latest_images[csi-attacher]="https://api.github.com/repos/kubernetes-csi/external-attacher/releases/latest"
latest_images[csi-provisioner]="https://api.github.com/repos/kubernetes-csi/external-provisioner/releases/latest"
latest_images[csi-snapshotter]="https://api.github.com/repos/kubernetes-csi/external-snapshotter/releases/latest"
latest_images[csi-resizer]="https://api.github.com/repos/kubernetes-csi/external-resizer/releases/latest"
latest_images[csi-node-driver-registrar]="https://api.github.com/repos/kubernetes-csi/node-driver-registrar/releases/latest"
latest_images[csi-external-health-monitor-controller]="https://api.github.com/repos/kubernetes-csi/external-health-monitor/releases/latest"
latest_images[otel-collector]="https://api.github.com/repos/open-telemetry/opentelemetry-collector/releases/latest"

retrieve_latest_tag() {
    curl -s $1 | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

update_versions() {
    version_file="config/csm-versions.yaml"

    for key in "${!latest_images[@]}"; do
        latest_tag=$(retrieve_latest_tag ${latest_images[$key]})
        if [ -z "$latest_tag" ]; then
            echo "Failed to retrieve latest image for $key. Try again later."
            exit 1
        fi


        if grep -q $key $version_file; then
            current_version=$(grep $key $version_file | cut -d ":" -f 2 | tr -d ' ')

            # Sanitize otel-collector version. Tag is different format than artifactory.
            if [ "$key" == "otel-collector" ]; then
                latest_tag=$(echo "$latest_tag" | sed 's/^v//')
            fi

            if [ "$current_version" == "$latest_tag" ]; then
                echo "$key already up to date"
                continue
            fi

            echo "Updating $key from $current_version to $latest_tag"
            sed -i "s/${key} : .*/${key} : ${latest_tag}/g" $version_file
        else
            echo "$key not found in version, adding."
            echo "$key : $latest_tag" >> $version_file
        fi
    done
}

update_versions
