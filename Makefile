all: build

CONTAINER_TOOL?=docker
IMAGE?="dell-common-installer"
DATA_COLLECTOR_IMAGE?="csm-data-collector"

VERSION?="v0.0.1"

build:
	go build -o dell-common-installer  .

run: build
	bash run.sh

docs:
	swag init

image: build
	$(CONTAINER_TOOL) build -t "$(IMAGE):$(VERSION)" .

image-push: image
	$(CONTAINER_TOOL) push "$(IMAGE):$(VERSION)"

data-collector:
	cd datacollectorapp && $(CONTAINER_TOOL) build -t "$(DATA_COLLECTOR_IMAGE):$(VERSION)" .

data-collector-push: data-collector
	$(CONTAINER_TOOL) push "$(DATA_COLLECTOR_IMAGE):$(VERSION)"

images: image-push data-collector-push

.PHONY: build run docs image image-push
