FROM ubuntu:20.10
RUN apt update
RUN apt install -y build-essential wget

WORKDIR /app

RUN wget -O kapp https://github.com/vmware-tanzu/carvel-kapp/releases/download/v0.36.0/kapp-linux-amd64 && chmod +x kapp

RUN wget -O kubectl https://storage.googleapis.com/kubernetes-release/release/v1.21.2/bin/linux/amd64/kubectl && chmod +x kubectl

COPY templates templates
COPY samplecerts samplecerts

COPY "dell-common-installer" .
COPY samplecerts samplecerts

CMD ["/app/dell-common-installer"]

EXPOSE 8080
