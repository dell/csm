# Common base image for all CSM images. When base image is upgraded, update the following 3 lines with URL, Version, and DEFAULT_BASEIMAGE variable.
# URL: https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58?architecture=amd64&image=65a8f97db7e4bede96526c22
# Version: ubi9/ubi-micro 9.4-1123
DEFAULT_BASEIMAGE="registry.redhat.io/ubi9/ubi@sha256:670f80d555117da5917dc4b7323926f47af5c6fade17cf55a174d37a0468a14a"
DEFAULT_GOIMAGE="golang:1.22"
