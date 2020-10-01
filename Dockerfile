FROM golang:1.15 as builder
WORKDIR /go/src/protoc-watch
COPY . .
RUN apt-get update && apt-get install -y unzip

# Download protoc
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip && \
    unzip protoc-3.6.1-linux-x86_64.zip -d protoc && \
    mv protoc/bin/* /usr/local/bin/ && \
    mv protoc/include/* /usr/local/include/

# Compile protoc-watch
RUN go get github.com/golang/protobuf/protoc-gen-go@v1.3.5 && \
    CGO_ENABLED=0 GOOS=linux go build


FROM alpine
COPY --from=builder /usr/local/include/google/* /usr/local/include/google/
COPY --from=builder /usr/local/bin/protoc /usr/local/bin/
COPY --from=builder /go/bin/protoc-gen-go /usr/local/bin/
COPY --from=builder /go/src/protoc-watch/protoc-watch /usr/local/bin/
RUN mkdir /home/protos
WORKDIR /home
