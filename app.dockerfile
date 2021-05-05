FROM golang:v2
WORKDIR /go
RUN go install github.com/go-delve/delve/cmd/dlv@latest &&\
    go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.2 &&\
    apt update &&\
    apt install unzip &&\
    cd bin &&\
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/protoc-3.15.6-linux-x86_64.zip &&\
    unzip protoc-3.15.6-linux-x86_64.zip bin/protoc &&\
    mv bin/protoc . &&\
    rm -r bin