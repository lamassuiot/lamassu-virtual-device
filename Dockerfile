FROM golang:1.18
WORKDIR /app
COPY . .
WORKDIR /app
ENV GOSUMDB=off
RUN CGO_ENABLED=0 go build -mod=vendor -o vdevice cmd/main.go

FROM ubuntu:22.04
WORKDIR /app
COPY --from=0 /app/vdevice /app
COPY config.tmpl /app/config.tmpl
RUN apt-get update && apt-get install wget curl -y
RUN curl -L https://github.com/a8m/envsubst/releases/download/v1.2.0/envsubst-`uname -s`-`uname -m` -o envsubst && chmod +x envsubst && mv envsubst /usr/local/bin
RUN wget https://www.amazontrust.com/repository/AmazonRootCA1.pem -o aws-iot-ca.crt
ENTRYPOINT cat /app/config.tmpl | envsubst > /app/config.json && ./vdevice
