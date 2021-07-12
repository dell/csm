#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Source common environment variables to gain access to
# HOST values, etc.
source $SCRIPT_DIR/../env.sh

# Ensure we can read in the local kube config.
KUBE_CONFIG=$HOME/.kube/config
if [[ ! -f $KUBE_CONFIG ]]
then
	echo no kube config at "$KUBE_CONFIG"
	exit 1
fi

# Ensure we can run jq.
if ! command -v jq > /dev/null 2>&1
then
	echo jq is required for this script
	exit 1
fi

URL="${SCHEME}://${HOST}:${PORT}"

# Create a new user (admin:admin)
DEPLOY_TOKEN=$( \
	curl -k -X POST "$URL/api/users" \
	-H "accept: application/json" \
	-H "Content-Type: application/json" \
	-d "{ \"user\": { \"password\": \"admin\", \"username\": \"admin\" }}" | jq -r '.user.token' \
)

echo "JWT token is: $DEPLOY_TOKEN"

# Create a new cluster entry.
curl -k -X POST "$URL/api/clusters" \
	-H "accept: application/json" \
	-H "Authorization: Bearer $DEPLOY_TOKEN" \
	-H "Content-Type: multipart/form-data" \
	-F "file=@$KUBE_CONFIG" -F "name=demo-k8s-cluster"

# Register a (powerflex) storage array.
curl -k -X POST "$URL/api/storage-arrays" \
	-H "accept: application/json" \
	-H "Authorization: Bearer $DEPLOY_TOKEN" \
	-H "Content-Type: application/json" \
	-d "{ \"storage-array\": { \"management_endpoint\": \"10.0.0.1\", \"password\": \"password\", \"storage_array_type\": \"powerflex\", \"unique_id\": \"id-1\", \"username\": \"user\" }}"

# Request deployment of a new application.
curl -k -X POST "$URL/api/applications" \
	-H "accept: application/json" \
	-H "Authorization: Bearer $DEPLOY_TOKEN" \
	-H "Content-Type: application/json" \
	-d "{ \"application\": { \"cluster_id\": 1, \"driver_type_id\": 0, \"module_types\": [ 0 ], \"name\": \"csi-powerflex-app\", \"storage_arrays\": [ 1 ] }}"
