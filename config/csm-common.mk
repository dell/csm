# Common base image for all CSM images. When base image is upgraded, update the following 3 lines with URL, Version, and DEFAULT_BASEIMAGE variable.
# URL: https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58?architecture=amd64&image=65a8f97db7e4bede96526c22
# Version: ubi9/ubi-micro 9.4-9
DEFAULT_BASEIMAGE="registry.redhat.io/ubi9/ubi-micro@sha256:979f04176cf3be5e0e70e246fd89184d2bcdde7762cf0b349e85dee5bc4cbb1f"
DEFAULT_GOIMAGE="golang:1.22"
