# Change this to your IP address if you want to access the server from elsewhere
export HOST="127.0.0.1"

# Change to the desired port
# If you are accessing it from a different host, make sure your firewall allows it 
export PORT="8080"

# Allowed values are https & http (recommended is https)
export SCHEME="https"

# Set it to a directory containing SSL certs
export CERT_DIR="samplecerts"

# SSL certificate (located in CERT_DIR)
export CERT_FILE="samplecert.crt"

# Key file (located in CERT_DIR)
export KEY_FILE="samplecert.key"

export DATA_COLLECTOR_IMAGE="localhost:5000/csm-data-collector:v0.0.1"

# kapp binary must be installed on host where service is running
export KAPP_BINARY="/root/kapp"
