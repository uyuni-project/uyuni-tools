# Tools to help using Uyuni as containers

**These tools are work in progress**

* `uyuniadm` used to help user administer uyuni servers on k8s and podman
* `uyunictl` used to help user managing Uyuni and SUSE Manager Servers mainly through its API

# Deployment rolling release

**NOTE:** This is rolling releases, meaning it can be broken at any time. Do not use it in production (yet!)

## For podman deployment
Requirement:
  - Opensuse Leap 15.4 or 15.5
  - podman installed

*Note that other distros with a recent podman installed could work, but for now the tool is not packaged for them in OBS.
So you would need to build it locally.*

Add uyuni-tool repository:
```
zypper ar https://download.opensuse.org/repositories/systemsmanagement:/Uyuni:/Master:/ContainerUtils/openSUSE_15.5/ uyuni-container-utils
```

Install `uyuniadm` package: `zypper in uyuniadm`

Run `uyuniadm` command to install uyuni server on podman:
```
uyuniadm install podman <MACHINE_FQDN>
```

If you built `uyuni-tools` on your machine, add the `--image registry.opensuse.org/systemsmanagement/uyuni/master/servercontainer/containers/uyuni/server` option to the install command.
This is not needed when using the package from OBS as it defaulting with this image at build time.

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
