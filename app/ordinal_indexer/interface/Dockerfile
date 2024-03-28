FROM golang:1.20 AS builder
ARG APP_RELATIVE_PATH
ARG CONFIG_FILE_NAME


COPY . /src
WORKDIR /src/app/${APP_RELATIVE_PATH}
RUN find ./configs -name "*.yaml" ! -name "${CONFIG_FILE_NAME}" -exec rm {} \;
RUN apt-get update \
    && apt-get install protobuf-compiler -y \
    && go mod download

RUN GOPROXY=https://goproxy.cn make build-docker

FROM debian:stable-slim
ARG APP_RELATIVE_PATH
ARG CONFIG_FILE_NAME
ENV CONFIG_FILE_NAME ${CONFIG_FILE_NAME}
RUN echo "Asia/Shanghai" > /etc/timezone && apt-get update && apt-get install ca-certificates -y && update-ca-certificates

COPY --from=builder /src/app/${APP_RELATIVE_PATH}/output /app
COPY --from=builder /go/bin/grpc-health-probe /usr/local/bin/grpc-health-probe

WORKDIR /app

EXPOSE 8080
EXPOSE 50051

CMD ["sh", "-c", "./bin/server -conf=./configs/$CONFIG_FILE_NAME"]
