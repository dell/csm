FROM ubuntu:20.10
RUN apt update
RUN apt install -y build-essential wget

WORKDIR /app

RUN wget -O kapp https://github.com/vmware-tanzu/carvel-kapp/releases/download/v0.36.0/kapp-linux-amd64 && chmod +x kapp

COPY templates templates
COPY samplecerts samplecerts

COPY "dell-common-installer" .

CMD ["/app/dell-common-installer"]

EXPOSE 8080
