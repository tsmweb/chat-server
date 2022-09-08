FROM golang:1.19

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apt-get update && \
    apt-get install build-essential protobuf-compiler -y && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

CMD ["tail", "-f", "/dev/null"]
