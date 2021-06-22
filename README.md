# csm
Dell EMC Container Storage Modules (CSM)

# How to build and deploy

Generate self-signed certificates using the following commands to add them to the `samplecerts` directory:
```
openssl req \
    -newkey rsa:4096 -nodes -sha256 -keyout samplecerts/csi_ca.key \
    -x509 -days 365 -out samplecerts/csi_ca.crt -subj '/'

openssl req \
    -newkey rsa:4096 -nodes -sha256 -keyout samplecerts/samplecert.key \
    -out samplecerts/samplecert.csr -subj '/'

openssl x509 -req -days 365 -in samplecerts/samplecert.csr -CA samplecerts/csi_ca.crt \
    -CAkey samplecerts/csi_ca.key -CAcreateserial -out samplecerts/samplecert.crt
```

Build images and push to local registry: `DATA_COLLECTOR_IMAGE=<registry-ip>:<registry-port>/csm-data-collector IMAGE=<registry-ip>:<registry-port>/dell-common-installer make images`

Edit deployment.yaml with location of images based on registry IP and port: `vi manifests/deployment.yaml`

Install dell-common-installer into kubernetes: `kubectl apply -f manifests/deployment.yaml`

If using Swagger, first access the REST API to accept the certificate in your browser (ex: `https://<k8s-node-ip>:31313/api/v1/users`)

After accepting the certificate, cd into the scripts directory: `cd scripts`

Edit the scripts/run-local-swagger.sh script with the IP to your kubernetes node.

Run run-local-swagger.sh: `./run-local-swagger.sh`

Open Swagger in your web browser using the IP address where the run-local-swagger.sh script is running: `http://<ip-address>:8080/swagger/`

