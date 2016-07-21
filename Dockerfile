FROM alpine:3.4

RUN apk --no-cache add  util-linux

COPY bin/linux/amd64/ /usr/bin/

