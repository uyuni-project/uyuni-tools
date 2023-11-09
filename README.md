<!--
SPDX-FileCopyrightText: 2023 SUSE LLC

SPDX-License-Identifier: Apache-2.0
-->

[![REUSE status](https://api.reuse.software/badge/git.fsfe.org/reuse/api)](https://api.reuse.software/info/git.fsfe.org/reuse/api)

# Tools to help using Uyuni as containers

**These tools are work in progress**

* `uyuniadm` used to help user administer Uyuni servers on K8s and Podman
* `uyunictl` used to help user managing Uyuni servers mainly through its API

# Deployment rolling release

**NOTE:** This is rolling releases, meaning it can be broken at any time. Do not use it in production (yet!)

## For Podman deployment
Requirement:
  - openSUSE Leap Micro 15.5
  - Podman installed

*Note that other distros with a recent Podman installed could work, but for now the tool is not packaged for them in OBS.
So you would need to build it locally.*

Add uyuni-tool repository:
```
zypper ar https://download.opensuse.org/repositories/systemsmanagement:/Uyuni:/Stable:/ContainerUtils/openSUSE_Leap_Micro_5.5/ uyuni-container-utils
```

Install `uyuniadm` package: `transactional-update pkg install uyuniadm`

Run `uyuniadm` command to install Uyuni server on Podman:
```
uyuniadm install podman
```

If you build `uyuni-tools` on your machine, add the `--image registry.opensuse.org/systemsmanagement/uyuni/stable/containers/uyuni/server` option to the install command.
This is not needed when using the package from OBS as it defaulting with this image at build time.

**NOTE**: rolling image url is: registry.opensuse.org/systemsmanagement/uyuni/master/containers/uyuni/server


Other sub-commands are also available. Explore the tool with the help command.

A tool named `uyunictl` is also available with useful commands.

## K3s deployment

For Look at a more details documentation at:

https://github.com/uyuni-project/uyuni/tree/master/containers/doc/server-kubernetes

# Development documentation

## Building

`go build -o ./bin ./...` will produce the binaries in the root folder with `0.0.0` as version.

To produce shell completion scripts for a given shell you can run:

- `./bin/uyuniadm completion <shell> > $COMPLETION_FILE` for uyuniadm
- `./bin/uyunictl error completion <shell> > $COMPLETION_FILE` for uyunictl

You'll then need to source the resulting script(s).

As an example, to enable bash completion for uyuniadm:

`./bin/uyuniadm completion bash > ./bin/completion`

`. ./bin/completion`

The supported shells are: bash, zsh and fish.

Alternatively, if you have `podman` installed you can run the `build.sh` script to build binaries compatible with any x86_64 linux.
The version will be computed from the last git tag and offset from it.

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

To disable features at build time pass the `-tags` parameter with the following values in a comma-separated list:

* `nok8s`: will disable Kubernetes support
