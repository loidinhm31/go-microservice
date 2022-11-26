## Modify module for import module in local environment
````
go mod edit -replace example.com/greetings=../greetings
````

## Connect Mongo Compass

````
mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&directConnection=true&ssl=false
````

## gRPC

### Install
Download pre-compiled binary for protoc

````
https://grpc.io/docs/protoc-installation/
````

Install the protocol compiler plugins for Go
````sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
````

Update bin path, example
````
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export GOBIN=$GOPATH/bin

export PATH=$PATH:$HOME/go/bin
````
Change <PATH_TO_PROTO_FILE> to use
````sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <PATH_TO_PROTO_FILE>
````

## Share volumes (file system) among nodes
Guide
````
https://phoenixnap.com/kb/sshfs
````

Tool
````
https://www.gluster.org/
````

## Kubernetes
minikube
````
https://minikube.sigs.k8s.io/docs/start/
````

kubernetes
````
https://kubernetes.io/docs/tasks/tools/
````

Configure Ingress TLS/SSL Certificates in Kubernetes
````
https://devopscube.com/configure-ingress-tls-kubernetes/
````