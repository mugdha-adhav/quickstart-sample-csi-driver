## Pre-requisites

### Setup
* Create a container image, containing your application binary.
* Create a `CSIDriver` resource.
* You might want to add RBAC resources according to your needs.
* Based on the [deployment strategy](#deployment-strategy), you may choose to create daemonset, statefulset, deployment, etc for your plugin.

### Deployment strategy
#### Simplest and easiest
* You may get started by just creating a daemonset.
* Thus, you don't need to package the node and controller components separately.
* You only need `node-driver-registrar` sidecar for running the plugin.

#### Containers required
1. Your plugin container
2. [NodeDriverRegistrar](https://kubernetes-csi.github.io/docs/node-driver-registrar.html)
3. [OPTIONAL] [Liveness probe](https://kubernetes-csi.github.io/docs/livenessprobe.html)

#### Volumes required
1. Driver socket file
    * Create a socket at the following path `/var/lib/kubelet/plugins/[CSIDriverName]/csi.sock` on the host and mount it in your pod at `/csi`.
    * This needs to be mounted on all containers.

2. Plugins registry
    * Mount `/var/lib/kubelet/plugins_registry` from the host at `/registration`.
    * Only mount in `node-driver-registrar` container.

3. Kubelet directory
    * Mount `/var/lib/kubelet` from the host at the same path on the pod.
    * Only mount in your plugin container.
    * [REQUIRED] [Bi-directional](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) mount propagation.


More details in [k8s docs](https://github.com/kubernetes/design-proposals-archive/blob/main/storage/container-storage-interface.md#kubelet-to-csi-driver-communication).
