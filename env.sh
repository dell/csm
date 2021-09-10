# Change this to your IP address if you want to access the server from elsewhere
export HOST="localhost"

# Change to the desired port
# If you are accessing it from a different host, make sure your firewall allows it
export PORT="31313"

# Allowed values are https & http (recommended is https)
export SCHEME="http"

# Set it to a directory containing SSL certs
export CERT_DIR="samplecerts"

# SSL certificate (located in CERT_DIR)
export CERT_FILE="samplecert.crt"

# Key file (located in CERT_DIR)
export KEY_FILE="samplecert.key"

export DATA_COLLECTOR_IMAGE="localhost:5000/csm-data-collector:v0.0.1"

# kapp binary must be installed on host where service is running
export KAPP_BINARY="/usr/local/bin/kapp"

#Username to access DB
export DB_USERNAME=root

#Password to access DB
export DB_PASSWORD=""

export DB_HOST=localhost

export DB_PORT=26257

export DB_SSL_ENABLED=false

export JWT_KEY=devjwtkey

export CIPHER_KEY=cipherkeyfordevelopmentpurposes!

# Your kubernetes node IP address
export API_SERVER_IP="localhost"
export API_SERVER_PORT="31313"
