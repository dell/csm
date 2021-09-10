<!--
Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
-->

# Getting Started Guide

## How to Build and Deploy CSM Installer

1. Build images and push to local registry: `DATA_COLLECTOR_IMAGE=<registry-ip>:<registry-port>/csm-data-collector IMAGE=<registry-ip>:<registry-port>/dell-common-installer make images`

**If securing the API service and database, following steps 2 to 4 to generate the certificates, or skip to step 5 to deploy without certificates**

2. From the `helm` directory, generate self-signed certificates using the following commands:

```
mkdir api-certs

openssl req \
    -newkey rsa:4096 -nodes -sha256 -keyout api-certs/ca.key \
    -x509 -days 365 -out api-certs/ca.crt -subj '/'

openssl req \
    -newkey rsa:4096 -nodes -sha256 -keyout api-certs/cert.key \
    -out api-certs/cert.csr -subj '/'

openssl x509 -req -days 365 -in api-certs/cert.csr -CA api-certs/ca.crt \
    -CAkey api-certs/ca.key -CAcreateserial -out api-certs/cert.crt
```

3. If required, download the `cockroach` binary used to generate certificates for the cockroach-db:
```
curl https://binaries.cockroachdb.com/cockroach-v21.1.8.linux-amd64.tgz | tar -xz && sudo cp -i cockroach-v21.1.8.linux-amd64/cockroach /usr/local/bin/
```

4. From the `helm` directory, generate the certificates required for the cockroach-db service:
```
mkdir db-certs

cockroach cert create-ca --certs-dir=db-certs --ca-key=db-certs/ca.key

cockroach cert create-node cockroachdb-0.cockroachdb.csm-installer.svc.cluster.local cockroachdb-public cockroachdb-0.cockroachdb --certs-dir=db-certs/ --ca-key=db-certs/ca.key

```
  In case multiple instances of cockroachdb are required add all nodes names while creating nodes on the certificates
```
cockroach cert create-node cockroachdb-0.cockroachdb.csm-installer.svc.cluster.local cockroachdb-1.cockroachdb.csm-installer.svc.cluster.local cockroachdb-2.cockroachdb.csm-installer.svc.cluster.local cockroachdb-public cockroachdb-0.cockroachdb cockroachdb-1.cockroachdb cockroachdb-2.cockroachdb --certs-dir=db-certs/ --ca-key=db-certs/ca.key
```
 
```
cockroach cert create-client root --certs-dir=db-certs/ --ca-key=db-certs/ca.key

cockroach cert list --certs-dir=db-certs/
```



5. Follow step `a` if certificates are being used or step `b` if certificates are not being used:

a) From the `helm` directory, install the helm chart, specifying the certificates generated in the previous steps:
```
helm install -n csm-installer --create-namespace \
   --set-file serviceCertificate=api-certs/cert.crt \
   --set-file servicePrivateKey=api-certs/cert.key \
   --set-file databaseCertificate=db-certs/node.crt \
   --set-file databasePrivateKey=db-certs/node.key \
   --set-file dbClientCertificate=db-certs/client.root.crt \
   --set-file dbClientPrivateKey=db-certs/client.root.key \
   --set-file caCrt=db-certs/ca.crt \
   csm-installer .
```
b) If not deploying with certificates, execute the following command:
```
helm install -n csm-installer --create-namespace \
   --set-string scheme=http \
   --set-string dbSSLEnabled="false" \
   csm-installer .
```