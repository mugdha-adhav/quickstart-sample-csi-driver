## Description
CSI driver for provisioning host-path volume.

## References
* [Developing CSI driver](https://kubernetes-csi.github.io/docs/developing.html).
* [CSI spec](https://github.com/container-storage-interface/spec/blob/master/spec.md).
* [Sample CSI driver](https://github.com/kubernetes-csi/csi-driver-host-path).
* [Blog](https://arslan.io/2018/06/21/how-to-write-a-container-storage-interface-csi-plugin/) on how to write CSI plugin.

## Development
### Local
Install [csc](https://github.com/rexray/gocsi/tree/master/csc) CLI tool for testing CSI RPCs locally.

Start running your application locally using -
```
go run cmd/main.go
```
Note the port mentioned for starting the server and `csc` is installed on your system.

Run below csc commands for testing your driver -
```
csc identity plugin-info --endpoint tcp://127.0.0.1:50052
```
In the above command `50052` is the port where the server is running locally.

Now you may test different identity, node and controller RPC methods by updating the above command.

## Build
### Build and run docker image
```
$ make build

$ make run
```

## Fix
[ ] Get multi stage docker build working.
