<!--
SPDX-FileCopyrightText: 2023 SUSE LLC

SPDX-License-Identifier: Apache-2.0
-->

The goal of this content is to set a high-level overview of each tool available.

For tools that depend on the backend we should explicitly specify which one we want to use. Backend can also be defined in the configuration file to be used by the user.

## Tools definition

In case one wants to add a new sub-command it should decide to in which tool it should be placed.

If the new sub-command needs access to the host OS of direct access to a running container then it should be added to MGRADM.

Commands in MGRCTL should use the API only. 

Any command to manage the proxy deployment must be placed in MGRPROXY.

MGRDEV is focused on utility commands to be used during the development process.


## MGRADM

**Goals and definition:**

Install, update, and maintain a containerized Uyuni Server. Commands placed here will have/need access to the container runtime environment and also to the HOST OS.

Any new command that needs direct access to the host OS or any running container must be added to these tools.

**Target Stakeholder:** Uyuni administrator
**Where to install:** System where Uyuni Server should be deployed
**Sub-commands Naming:** verb -> backend

## MGRCTL

**Goals and definition:**
Helper tool for day-to-day operations and integration with other tools.
Sub-commands in this tool should use the API calls (although case-by-case exceptions can be considered if there are valid reasons).

**Target Stakeholder:** Uyuni operators
**Where to install:** System where the Uyuni Server is deployed, or in the operator machine (supporting the same Operating Systems we already support for `spacecmd`).
**Sub-commands Naming:** subcommand -> verb

## MGRDEV

**Goals and definition:**
Utility commands to be used during development process. This tool can have commands that run remotely on the host OS or on running containers. These commands can use SSH and podman-socket.
Examples of sub-commands are `cp` and `exec`.

**Target Stakeholder:** Uyuni Developers
**Where to install:** Any machine that needs remote access to running containers.
**Sub-commands Naming:** verb -> backend


## MGRPROXY

**Goals and definition:**
Install and manage a containerized Uyuni Proxy. This new command is a proposal to solve the problem of managing the Proxy using the same tool that manages the server, and how that can lead to confusion and errors.

**Target Stakeholder:** Uyuni administrator
**Where to install:** System where the Uyuni Proxy should be deployed
**Sub-commands Naming:** verb -> backend

This command is to be developed in a later stage since it would be better to redefine how we deploy containerized proxy and follow the same approach we have provided in the server.