FROM golang:1.6

ENV GO_PACKAGE=github.com/wikiwi/kube-volume-freezer

RUN mkdir -p /go/src/${GO_PACKAGE}
COPY . /go/src/${GO_PACKAGE}

WORKDIR /go/src/${GO_PACKAGE}

RUN make bootstrap build && \
    mv /go/src/${GO_PACKAGE}/bin/* /go/bin/ && \
    rm -rf /go/pkg

