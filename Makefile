all: build

CONTAINER_TOOL?=docker
IMAGE?=$(IMAGE_REPO)dell-csm-installer
DATA_COLLECTOR_IMAGE?=$(IMAGE_REPO)csm-data-collector
BASE_IMG?=$(BASE_IMAGE)

VERSION?="v0.0.1"

build:
	go mod vendor
	go build -o dell-common-installer  .

run: build
	bash run.sh

generate:
	go generate ./...

test:
	go test -cover -race -count 1 -timeout 5m ./...

docs:
	swag init

image: build
ifeq ($(BASE_IMG),)
	$(CONTAINER_TOOL) build --build-arg BASE_IMG="ubuntu:20.10" -t "$(IMAGE):$(VERSION)" .
else
	$(CONTAINER_TOOL) build --build-arg BASE_IMG="$(BASE_IMG)" -t "$(IMAGE):$(VERSION)" .
endif

image-replication: build
	if [ -d "../dell-csi-replicator" ]; then \
        echo "Replication directory found building the repctl binary"; \
		cp -rf ../dell-csi-replicator .; \
	else \
		echo "Replication directory not found please clone the replication project to ../dell-csi-replicator and restart the build"; \
		exit 1; \
    fi
	if [ -d "../dell-csi-extensions" ]; then \
        echo "Extensions directory found building the repctl binary"; \
		cp -rf ../dell-csi-extensions .; \
	else \
		echo "Extensions directory not found please clone the extensions project to ../dell-csi-extensions and restart the build"; \
		exit 1; \
    fi
ifeq ($(BASE_IMG),)
	$(CONTAINER_TOOL) build -f Dockerfile.dev --build-arg BASE_IMG="golang:1.16.7" -t "$(IMAGE):$(VERSION)" .
else
	$(CONTAINER_TOOL) build -f Dockerfile.dev --build-arg BASE_IMG="$(BASE_IMG)" -t "$(IMAGE):$(VERSION)" .
endif
	rm -rf dell-csi-replicator
	rm -rf dell-csi-extensions

image-push: image
	$(CONTAINER_TOOL) push "$(IMAGE):$(VERSION)"

image-push-replication: image-replication
	$(CONTAINER_TOOL) push "$(IMAGE):$(VERSION)"

# Targets to start the Cockroach DB locally in windows and unix systems for development and unit tests
win-install-local-db:
	wget https://binaries.cockroachdb.com/cockroach-v21.1.8.windows-6.2-amd64.zip -O temp.zip
	unzip temp.zip
	rm -rf temp.zip
	mv -f cockroach-v21.1.8.windows-6.2-amd64/ .cockroach

unix-install-local-db:
	wget https://binaries.cockroachdb.com/cockroach-v21.1.8.linux-amd64.tgz -O temp.tgz --no-check-certificate
	tar zxvf temp.tgz
	rm -rf temp.tgz
	mv -f cockroach-v21.1.8.linux-amd64/ .cockroach

start-local-db:
	./.cockroach/cockroach start-single-node --insecure

# Stop the local DB with CTRL+C before uninstall the files
uninstall-local-db:
	rm -rf .cockroach
	rm -rf cockroach-data

data-collector:
	cd datacollectorapp && $(CONTAINER_TOOL) build -t "$(DATA_COLLECTOR_IMAGE):$(VERSION)" .

data-collector-push: data-collector
	$(CONTAINER_TOOL) push "$(DATA_COLLECTOR_IMAGE):$(VERSION)"

images: image-push data-collector-push

.PHONY: build run docs image image-push
