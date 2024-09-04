# Common base image for all CSM images. When base image is upgraded, do not update anything as it is automated.
# URL: https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58?architecture=amd64
# Version: ubi9/ubi-micro 9.4-15
DEFAULT_BASEIMAGE="registry.access.redhat.com/ubi9/ubi-micro@sha256:7f376b75faf8ea546f28f8529c37d24adcde33dca4103f4897ae19a43d58192b"
DEFAULT_GOIMAGE="golang:1.22"