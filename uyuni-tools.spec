#
# spec file for package uyuni-tools
#
# Copyright (c) 2023 SUSE LLC
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative.

# Please submit bugfixes or comments via https://bugs.opensuse.org/
#

%global provider        github
%global provider_tld    com
%global org             uyuni-project
%global project         uyuni-tools
%global provider_prefix %{provider}.%{provider_tld}/%{org}/%{project}
%global productname     Uyuni

%global image           registry.opensuse.org/uyuni/server
%global chart           oci://registry.opensuse.org/uyuni/server-helm

%if 0%{?suse_version} >= 1600 || 0%{?sle_version} >= 150400 || 0%{?rhel} >= 8 || 0%{?fedora} >= 37 || 0%{?debian} >= 12 || 0%{?ubuntu} >= 2004
%define adm_build    1
%else
%define adm_build    0
%endif


Name:           %{project}
Version:        0.0.1
Release:        0
Summary:        Tools for managing %{productname} container
License:        Apache-2.0
Group:          System/Management
URL:            https://%{provider_prefix}
Source0:        %{name}-%{version}.tar.gz
Source1:        vendor.tar.gz
BuildRequires:  coreutils
# Get the proper Go version on different distros
%if 0%{?suse_version}
BuildRequires:  golang(API) >= 1.20
%endif
%if 0%{?ubuntu}
%define go_version      1.20
BuildRequires:  golang-%{go_version}
%endif
%if 0%{?debian}
BuildRequires:  golang >= 1.20
%endif
%if 0%{?fedora} || 0%{?rhel}
BuildRequires:  golang >= 1.19
%endif

%description
Tools for managing uyuni container.

%if %{adm_build}
%package -n uyuniadm
Summary:      Command line tool to install and update %{productname}

%description -n uyuniadm
uyuniadm is a convenient tool to install and update %{productname} components as containers running
either on Podman or a Kubernetes cluster.
%endif

%package -n uyunictl
Summary:      Command line tool to perform day-to-day operations on %{productname}

%description -n uyunictl
uyunictl is a tool helping with dayly tasks on %{productname} components running as containers
either on Podman or a Kubernetes cluster.


%prep
%autosetup
tar -zxf %{SOURCE1}


%build
export GOFLAGS=-mod=vendor
mkdir -p bin
ADM_PATH="%{provider_prefix}/uyuniadm/shared/utils"
UTILS_PATH="%{provider_prefix}/shared/utils"

tag=%{!?_default_tag:latest}
%if "%{?_default_tag}" != ""
    tag='%{_default_tag}'
%endif

image=%{image}
%if "%{?_default_image}" != ""
  image='%{_default_image}'
%endif

chart=%{chart}
%if "%{?_default_chart}" != ""
  chart='%{_default_chart}'
%endif

go_tags=""
%if "%{?_uyuni_tools_tags}" != ""
  go_tags="-tags %{_uyuni_tools_tags}"
%endif

go_path=
%if 0%{?ubuntu}
  go_path=/usr/lib/go-%{go_version}/bin/
%else
  %if "%{?_go_bin}" != ""
    go_path='%{_go_bin}/'
  %endif
%endif

GOLD_FLAGS="-X ${ADM_PATH}.DefaultImage=${image} -X ${ADM_PATH}.DefaultTag=${tag} -X ${ADM_PATH}.DefaultChart=${chart}"
GOLD_FLAGS="${GOLD_FLAGS} -X ${UTILS_PATH}.Version=%{version}"

# Workaround for rpm on Fedora and EL clones not able to handle go's compressed debug symbols
# Found compressed .debug_aranges section, not attempting dwz compression
%if 0%{?rhel} >= 8 || 0%{?fedora} >= 38
GOLD_FLAGS="-compressdwarf=false ${GOLD_FLAGS}"
%endif

# Workaround for missing build-id on Fedora
# error: Missing build-id in [...]
%if 0%{?fedora} >= 38
GOLD_FLAGS="-B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') ${GOLD_FLAGS}"
%endif

${go_path}go build ${go_tags} -ldflags "${GOLD_FLAGS}" -o ./bin ./...

%if ! %{adm_build}
rm ./bin/uyuniadm
%endif

%install
install -m 0755 -vd %{buildroot}%{_bindir}
install -m 0755 -vp ./bin/* %{buildroot}%{_bindir}/

%if %{adm_build}
%files -n uyuniadm
%defattr(-,root,root)
%doc README.md
%license LICENSE
%{_bindir}/uyuniadm
%endif

%files -n uyunictl
%defattr(-,root,root)
%doc README.md
%license LICENSE
%{_bindir}/uyunictl

%changelog
