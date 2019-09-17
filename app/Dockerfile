ARG GO_VERSION=1.13
FROM golang:${GO_VERSION}-alpine
LABEL Description="Testing stream movie" Vendor="Alex Peréz Ortuño" Version="1.0" maintainer="Alex Peréz <alexperezortuno@gmail.com>"

RUN apk update && apk upgrade && apk add --no-cache bash git openssh ca-certificates

WORKDIR /go/src
COPY src/api /go/src/api/
COPY src/main.go /go/src/
RUN go get -t github.com/gorilla/mux
RUN go get -t gopkg.in/natefinch/lumberjack.v2
RUN go build -o stream_api

EXPOSE 3000

ENTRYPOINT /go/src/stream_api