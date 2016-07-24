FROM alpine:3.4

ARG BUILD_DATE
ARG BUILD_URL
ARG VCS_REF
ARG VCS_VERSION
ARG VCS_MESSAGE

LABEL org.label-schema.build-date=${BUILD_DATE} \
      org.label-schema.vcs-ref=${VCS_REF} \
      org.label-schema.vcs-version=${VCS_VERSION} \
      org.label-schema.vcs-url="https://github.com/wikiwi/kube-volume-freezer" \
      org.label-schema.vendor=wikiwi.io \
      org.label-schema.name=kube-volume-freezer \
      io.wikiwi.build-url=${BUILD_URL} \
      io.wikiwi.license=MIT \
      io.wikiwi.vcs-msg=${VCS_MESSAGE}

RUN apk --no-cache add util-linux

COPY bin/linux/amd64/ /usr/bin/

