# Common base image for all CSM images. When base image is upgraded, update the following 3 lines with URL, Version, and DEFAULT_BASEIMAGE variable.
# URL: https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58?architecture=amd64&image=65a8f97db7e4bede96526c22
# Version: ubi9/ubi-micro 9.3-13
DEFAULT_BASEIMAGE="registry.access.redhat.com/ubi9/ubi-micro@sha256:d72202acf3073b61cb407e86395935b7bac5b93b16071d2b40b9fb485db2135d"
DEFAULT_GOIMAGE="golang:1.22"
