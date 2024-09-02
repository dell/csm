# Common base image for all CSM images. When base image is upgraded, do not update anything as it is automated.
# URL: https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58?architecture=amd64
# Version: ubi9/ubi-micro 9.4-13
DEFAULT_BASEIMAGE="registry.access.redhat.com/ubi9/ubi-micro@sha256:9dbba858e5c8821fbe1a36c376ba23b83ba00f100126f2073baa32df2c8e183a"
DEFAULT_GOIMAGE="golang:1.22"