
# Download swagger.json and replace host with correct IP
IP=<ip-address-here>:31313 && curl --insecure https://$IP/swagger/doc.json > tmp.json && jq ".host=\"$IP\"" tmp.json > swagger.json && rm tmp.json

# Serve UI with container
docker run --network=host -e BASE_URL=/swagger -e SWAGGER_JSON=/swagger/swagger.json -v $(pwd):/swagger swaggerapi/swagger-ui
