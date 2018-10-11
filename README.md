# protoc-watch

Watch `.proto` files in a folder and auto compile them on creation or change.

### Install
```sh
go get github.com/kvartborg/protoc-watch
```

### Use in docker-compose
```yml
version: 3
services:
  dev-protoc-watch:
    image: kvartborg/protoc-watch
    command: protoc-watch --go_out=. ./protos
    volumes:
      - src/protos:/home/protos
```
