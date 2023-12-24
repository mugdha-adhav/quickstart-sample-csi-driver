## Installation
Install [csc](https://github.com/rexray/gocsi/tree/master/csc) CLI tool for testing CSI RPCs locally.

## Pre-requisite
* Your application should start a gRPC server on a specific port.
* Refer [CSI spec](https://github.com/container-storage-interface/spec/blob/master/spec.md) for getting details on functions to be implemented. 

## Start running your app
For running your application locally without any setup, you may directly run your binary/code like -
```sh
go run cmd/main.go
```

**Note**: The port number on which the gRPC server is running is needed while running csc commands.

## Testing your driver
Sample for testing GetPluginInfo() identity function -
```sh
csc identity plugin-info --endpoint tcp://127.0.0.1:50052
```
**Note**: In the above command `50052` is the port where the server is running for this application.
You may replace it the port number noted in [Start running your app](#start-running-your-app) section

Next, test different identity, node and controller RPC methods by updating the above command.
