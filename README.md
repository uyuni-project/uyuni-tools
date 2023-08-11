# Tools to help using Uyuni as containers

**These tools are work in progress**

* `uyuniadm` used to help user administer uyuni servers on k8s and podma
* `uyunictl` used to help user managing Uyuni and SUSE Manager Servers mainly through its API

## Building

`go build -o ./bin ./...` will produce the binaries in the root folder.
Alternatively, if you have `podman` installed you can run the `build.sh` script to build binaries compatible with SLE 15 SP4 or openSUSE Leap 15.4.


## Uyuniadm usage

Available Commands:
  * **help**: Help about any command
  * **install**: install a new server from scratch
  * **migrate**: migrate a remote server to containers
  * **uninstall**: uninstall a server

For more information about flags `uyuniadm --help`

### Uyuniadm install

Install a new server from scratch

The install command assumes the following:
  * podman or kubectl is installed locally
  * if kubectl is installed, a working kubeconfig should be set to connect to the cluster to deploy to

When installing on kubernetes, the helm values file will be overridden with the values from the uyuniadm parameters or configuration.

NOTE: for now installing on a remote cluster or podman is not supported!

```
Usage:
  uyuniadm install [fqdn] [flags]
```

For more information about flags `uyuniadm install --help`

### Uyuniadm migrate

Migrate a remote server to containers

This migration command assumes a few things:
  * the SSH configuration for the source server is complete, including user and
    all needed options to connect to the machine,
  * an SSH agent is started and the key to use to connect to the server is added to it,
  * podman or kubectl is installed locally
  * if kubectl is installed, a working kubeconfig should be set to connect to the cluster to deploy to

NOTE: for now installing on a remote cluster or podman is not supported yet!

```
Usage:
  uyuniadm migrate [source server FQDN] [flags]
```

#### SSH Configuration Example
1. In the destination server, add to `~/.ssh/config` :
   ```
   Host YOUR_HOST
    Hostname SOURCE_HOSTNAME
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
    Port 22
    User SOURCE_USER
    IdentitiesOnly yes
    ```
2. If you already have a key, run:

    ```
    ssh-copy-id YOUR_HOST
    ```
    If not, run `ssh-keygen` to generate it.
3. If the `SOURCE_USER` user is not root, it should be able to run `rsync`. It can be done by adding to `/etc/sudoers`:
    ```
    add to sudoers file
    SOURCE_USER ALL=(ALL) NOPASSWD:/usr/bin/rsync
    ```
4. To provide a ssh agent with key, in the destination server:
    ```
    eval `ssh-agent`
    ssh-add $KEY_PATH
    ```
### Uyuniadm uninstall

Uninstall a server

For more information about flags `uyuniadm uninstall --help`


## Uyunictl usage
Available Commands:
  * **cp**: copy files to and from the containers
  * **exec**: execute commands inside the uyuni containers using 'sh -c'
  * **help**: Help about any command

### Uyunictl cp

Takes a source and destination parameters. One of them can be prefixed with 'server:' to indicate the path is within the server pod.

```
Usage:
  uyunictl cp [path/to/source.file] [path/to/destination.file] [flags]
```
For more information about flags `uyunictl cp --help`

### Uyunictl exec

Execute commands inside the uyuni containers using 'sh -c'

```
Usage:
  uyunictl exec '[command-to-run --with-args]' [flags]
```
For more information about flags `uyunictl exec --help`
