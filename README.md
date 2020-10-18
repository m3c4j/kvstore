# KVStore
Learning a variety of tree structures by creating a k-v store database.



### Build





### Prepare

This project using [grpc](https://www.grpc.io/docs/languages/go/quickstart/) to communicate from kvctl to kvd, we need protobuf and protoc-gen-go installed.

```bash
# on debian like linux system
sudo apt install -y protobuf-compiler

export GO111MODULE=on  # Enable module mode
go get github.com/golang/protobuf/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0

export PATH="$PATH:$(go env GOPATH)/bin"

```

##### grpc

we are using gogo protobuf library

https://github.com/gogo/protobuf

```bash
go get github.com/gogo/protobuf/protoc-gen-gofast

# we using this repo's proto files
git clone https://github.com/googleapis/googleapis

# generate pb.go file
protoc -I=. -I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.1 -I/mnt/e/Develop/Projects/GitHub/googleapis --gofast_out=plugins=grpc:. rpc.proto

# same as
protoc -I . \
	-I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.1 \
	-I/mnt/e/Develop/Projects/GitHub/googleapis \
	--grpc-gateway_out . \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	rpc.proto

# https://github.com/googleapis/go-genproto
go get google.golang.org/genproto/...

```









