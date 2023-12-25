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

1. Plugins registry
    * Mount `/var/lib/kubelet/plugins_registry` from the host at `/registration`.
    * Only mount in `node-driver-registrar` container.

1. Kubelet pods directory
    * Mount `/var/lib/kubelet/pods` from the host at the same path on the pod.
    * Only mount in your plugin container.
    * [REQUIRED] [Bi-directional](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) mount propagation.

1. Kubelet plugins directory
    * Mount `/var/lib/kubelet/plugins` from the host at the same path on the pod.
    * Only mount in your plugin container.
    * [REQUIRED] [Bi-directional](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation) mount propagation.


More details in [k8s docs](https://github.com/kubernetes/design-proposals-archive/blob/main/storage/container-storage-interface.md#kubelet-to-csi-driver-communication).

### Issues faced
#### Mounting kubelet directory instead of child pod and plugins directory
In the design proposals archive doc for CSI, it's mentioned to mount only the kubelet directory, instead of pods and plugins kubelet directories separately.

**What went wrong**
I tried mounting `/var/lib/kubelet` directory from the host onto the plugins container, but it didn't mount the host-path onto the target-path.

**Debugging steps**
While trying to create the target-directory in `NodePublishVolume` function, it was failing with error path doesn't exist.

To workaround this, I tried creating the target-directory with `MkdirAll()` function, which would create the entire path even if it didn't exist, but it still didn't mount the expected target-directory.

Also checked the mount paths on the host using `mount` command and listing contents in `/proc/mounts` directory but this didn't help.

Finally while getting list of mounted filesystems from `/etc/mtab` which stores the mount and umount commands that were executed, I found that the mount was hapenning on `/var/lib/kubelet/<pod-id>` instead of `/var/lib/kubelet/pod/<pod-id>`.

**Fix:**
After mounting the `/var/lib/kubelet/pods` and `/var/lib/kubelet/plugins` directories separately on plugin container, the mounting of host-path on target-path worked.
