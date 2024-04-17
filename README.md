<!--
SPDX-FileCopyrightText: 2023 SUSE LLC

SPDX-License-Identifier: Apache-2.0
-->

[![REUSE status](https://api.reuse.software/badge/git.fsfe.org/reuse/api)](https://api.reuse.software/info/git.fsfe.org/reuse/api)

# Tools to help using Uyuni as containers

**These tools are work in progress**

* `mgradm` used to help administer Uyuni servers on K8s and Podman
* `mgrctl` used to help managing Uyuni servers mainly through its API
* `mgrpxy` used to help managing Uyuni proxies

# Deployment rolling release

## For Podman deployment
Requirement:
  - openSUSE Leap Micro 15.5
  - Podman installed

*Note that other distros with a recent Podman installed could work but they have not been tested.
Please report issues if any arises on those distributions.*

Add uyuni-tool repository:
```
zypper ar https://download.opensuse.org/repositories/systemsmanagement:/Uyuni:/Stable:/ContainerUtils/openSUSE_Leap_Micro_5.5/ uyuni-container-utils
```

Install `mgradm` package: `transactional-update pkg install mgradm`

Run `mgradm` command to install Uyuni server on Podman:
```
mgradm install podman
```

If you build `uyuni-tools` on your machine, add the `--image registry.opensuse.org/systemsmanagement/uyuni/stable/containers/uyuni/server` option to the install command.
This is not needed when using the package from OBS as it defaulting with this image at build time.

**NOTE**: rolling image url is: registry.opensuse.org/systemsmanagement/uyuni/master/containers/uyuni/server


Other sub-commands are also available. Explore the tool with the help command.

A tool named `mgrctl` is also available with useful commands.

## K3s deployment

For Look at a more details documentation at:

https://github.com/uyuni-project/uyuni/tree/master/containers/doc/server-kubernetes

# Development documentation

## Building

`go build -o ./bin ./...` will produce the binaries in the root folder with `0.0.0` as version.

To produce shell completion scripts for a given shell you can run:

- `./bin/mgradm completion <shell> > $COMPLETION_FILE` for mgradm
- `./bin/mgrctl error completion <shell> > $COMPLETION_FILE` for mgrctl

You'll then need to source the resulting script(s).

As an example, to enable bash completion for mgradm:

`./bin/mgradm completion bash > ./bin/completion`

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

## Localization

### Developer tricks

For Localization the project uses `gettext`.
There are a few rules to follow to make strings localizable:

Add the following import in the go file and then wrap all the strings that could be localized in the `L()` function.

```
. "github.com/uyuni-project/uyuni-tools/shared/l10n"
```

**Global variables and constants are evaluated before running the main function and thus do not take the locale into account.**
Move them in a function to work around this issue.

### Generating the POT files

In order to extract the strings from the code run the `extract_strings` script.
One POT file for each tool and one for the `shared` folder will be generated in the `locale` directory.

### Translating

The translation files should be named after the target language next to the corresponding PO file.
The `.mo` files should not be committed in the source tree as they are build results.
Those are generated using the `locale/build.sh` script.
