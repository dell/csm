# Common base image for all CSM images. When base image is upgraded, update the following 3 lines with URL, Version, and DEFAULT_BASEIMAGE variable.
# URL: https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58?architecture=amd64&image=65a8f97db7e4bede96526c22
# Version: ubi9/ubi-micro 9.4-9
DEFAULT_BASEIMAGE="registry.access.redhat.com/ubi9-micro@sha256:2044e2ca8e258d00332f40532db9f55fb3d0bfd77ecc84c4aa4c1b7af3626ffb"
DEFAULT_GOIMAGE="golang:1.22"
