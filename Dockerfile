ARG BASE_IMG
FROM $BASE_IMG AS builder

WORKDIR /app

RUN apt update
RUN apt install -y build-essential wget

RUN wget -O kapp https://github.com/vmware-tanzu/carvel-kapp/releases/download/v0.40.0/kapp-linux-amd64 && chmod +x kapp

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4-208

WORKDIR /app

COPY templates templates
COPY --from=builder /app/kapp .
COPY "dell-common-installer" .

CMD ["/app/dell-common-installer"]

EXPOSE 8080
