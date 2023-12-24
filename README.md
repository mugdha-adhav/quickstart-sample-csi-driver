## Description
CSI driver for mounting a volume provided by the driver inside a pod.

## Getting started
1. Decide the name for your CSI driver. It should be a standard/sanitized name that doesn't contain dangerous characters.
1. Create a gRPC server which listens on a specific port.
1. Start implementing the methods defined in [CSI spec](https://github.com/container-storage-interface/spec/blob/master/spec.md), initial implementation can just be an empty function definition.
1. You may verify whether you have added required implementations as expected using `csc` tool. More details on how to test using csc can be found in [csc.md](./docs/csc.md).
1. NodeGetInfo() and GetPluginInfo() functions need to be have valid implementations for the driver to get registered.
    * The parameter values can be passed in from config or hard-coded (not recommended).
    * GetPluginInfo() configures the driver name and version.
    * NodeGetInfo() returns the node ID in response.
1. Register the Node and Identity components using csi library like `csi.RegisterIdentityServer(serverOptions, driverOptions)`.
1. Logic for mount and unmount volume operations goes in NodePublishVolume() and NodeUnpublishVolume() functions.
1. Docs for [Configuring and Deploying](./docs/configure-and-deploy.md) your application explains how to setup and get your app running.
1. After the deployment completes, you may verify if your driver is functional by following [verification docs](./docs/verify.md)

## ToDo
### Create a directory on the host and mount it.
- [x] Setup gRPC server
- [x] Expose Controller, Node and Identity Service RPCs
- [x] Identity Service implemented
- [x] Node Service implemented
    - [x] `NodePublishVolume` implemented
        - Volume corresponding to the `volume_id` published(mounted) at specified `target_path`
    - [x] `NodeUnpublishVolume` implemented
        - This RPC is a reverse operation of `NodePublishVolume`.
    - [x] `NodeGetCapabilities` implemented

## Development
### Build and push image
#### Build docker image
```sh
make build TAG="your-tag"
```

#### Push docker image to registry
```sh
make push TAG="your-tag"
```

#### Load docker image in kind cluster
```sh
make load TAG="your-tag"
```

### Run/Deploy CSI driver
#### Run the application in docker container
```sh
make run TAG="your-tag"
```

#### Deploy manifests on kind cluster
```sh
make deploy-kind TAG="your-tag"
```

### Cleanup
### Delete all resources from kind cluster
```sh
make remove-kind TAG="your-tag"
```

## References
* [Developing CSI driver](https://kubernetes-csi.github.io/docs/developing.html).
* [Design proposal](https://github.com/kubernetes/design-proposals-archive/blob/main/storage/container-storage-interface.md)
* [Sample CSI driver](https://github.com/kubernetes-csi/csi-driver-host-path).
* [Blog](https://arslan.io/2018/06/21/how-to-write-a-container-storage-interface-csi-plugin/) on how to write CSI plugin.
