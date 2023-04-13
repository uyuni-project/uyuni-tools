# Tools to help using Uyuni as containers

**These tools are work in progress**

* `uyunictl` aims at providing utility functions for day-to-day operations against a containerized Uyuni server
* `uyuniadm` aims at helping deployment and setup of a containerized Uyuni server

## Building

`go build -o ./bin ./...` will produce the binaries in the root folder.
Alternatively, if you have `podman` installed you can run the `build.sh` script to build binaries compatible with SLE 15 SP4 or openSUSE Leap 15.4.
