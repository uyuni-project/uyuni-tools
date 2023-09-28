# Tools to help using Uyuni as containers

**These tools are work in progress**

* `uyuniadm` used to help user administer uyuni servers on k8s and podman
* `uyunictl` used to help user managing Uyuni and SUSE Manager Servers mainly through its API

# Deployment rolling release

**NOTE:** This is rolling releases, meaning they can break at any time. Do not use it in product (yet!)

## For podment deployment
Requirement:
  - Opensuse Leap 15.4 or 15.5
  - podman installed

Add uyuni-tool repository:
```
zypper ar https://download.opensuse.org/repositories/systemsmanagement:/Uyuni:/Master:/ServerContainer/openSUSE_Leap_15.4/ uyuni-tools
```

Install uyuni-tool package: `zypper in uyuni-tools`

Run `uyuniadm` command to install uyuni server on podman:
```
uyuniadm install podman --image registry.opensuse.org/systemsmanagement/uyuni/master/servercontainer/containers/uyuni/server <MACHINE_FQDN>
```

**NOTE**: rolling image url is: registry.opensuse.org/systemsmanagement/uyuni/master/servercontainer/containers/uyuni/server


Other sub-cammands are also available. Explore the tool with the help command.

A tool named `uyunictl` is also available with usefull commands.

## k3s deployment

For Look at a more details documentation at:

https://github.com/uyuni-project/uyuni/tree/server-container/containers/doc/server-kubernetes

# Technical documentation

## Building

`go build -o ./bin ./...` will produce the binaries in the root folder.
Alternatively, if you have `podman` installed you can run the `build.sh` script to build binaries compatible with SLE 15 SP4 or openSUSE Leap 15.4.

### Building in Open Build Service

In order to adjust the image, tag and chart to the project the package is built in, add the following at the end of the project configuration:

```
Macros:
%_default_tag yourtag
%_default_image theregistry.org/path/to/the/server
%_default_chart oci://theregistry.org/path/to/the/chart
:Macros
```

### Disabling features at build time

To disable features at build time pass the `-tags` paramter with the following values in a comma-separated list:

* `nok8s`: will disable kubernetes support

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
For more information about flags `uyuniadm migrate --help`

#### SSH Configuration Example
1. In the destination server, add to `~/.ssh/config` :
   ```
   Host SOURCE_HOSTNAME
    Hostname SOURCE_HOSTNAME
    StrictHostKeyChecking no
    Port 22
    User SOURCE_USER
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
```
Usage:
  uyuniadm uninstall [flags]
```

For more information about flags `uyuniadm uninstall --help`


## Uyunictl usage
Available Commands:
  * **cp**: copy files to and from the containers
  * **exec**: execute commands inside the uyuni containers using 'sh -c'
  * **help**: Help about any command

Using `uyunictl` to access a remote cluster requires `kubectl` to be configured to connect to it before hand.

In order to connect to a remote `podman`, ensure the `podman.socket` systemd unit is active on the server by running `systemctl enable --now podman.socket`.
Then configure the Podman connection on the client machine:

```
podman system connection add <name> ssh://root@<host.fqdn>
```

Then export `CONTAINER_CONNECTION=<name>` before running `uyunictl`.
Note that passing `--identity <file>` may be needed to tell SSH which key to use to connect to the podman host.


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

## Configuration File Example
All the commands can accept flags or yaml configuration file (using the option `-c`). This is an example of configuration file:
```
db:
  password: YOUR_DB_PASSWORD
cert:
  password: YOUR_DB_PASSWORD
scc:
  user: YOUR_SCC_USER
  password: YOUR_SCC_PASSWORD
email: YOUR_MAIL
emailFrom: YOUR_MAIL
image: YOUR_IMAGE_REGISTRY

helm:
  uyuni:
    chart: oci://OCI_REGISTRY
    values: /root/chart-values.yaml
podman:
  arg:
    - -p
    - 8000:8000
    - -p
    - 8001:8001
    - ""
```


# Podman Deployment Example 
 
Requirements for the Host OS:
  - Use Leap 15.5
  - Have podman installed
  - Have a valid FQDN for the machine

Create a file "/root/uyuniadm.yaml" with the following content:

```
db:
  password: spacewalk
cert:
  password: spacewalk
image: registry.opensuse.org/systemsmanagement/uyuni/master/servercontainer/containers/uyuni/server
```

Then run `uyuniadm install --config /root/uyuniadm.yaml MACHINE-FQDN`
