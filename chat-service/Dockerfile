FROM golang:1.16

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apt-get update && \
    apt-get install build-essential protobuf-compiler -y && \
    go get google.golang.org/protobuf/cmd/protoc-gen-go && \
    go get github.com/google/wire/cmd/wire

CMD ["tail", "-f", "/dev/null"]