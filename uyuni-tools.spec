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

Name:           %{project}
Version:        0.0.1
Release:        0
Summary:        Tools for managing uyuni container
License:        Apache-2.0
Group:          System/Management
URL:            https://%{provider_prefix}
Source0:        %{name}-%{version}.tar.gz
Source1:        vendor.tar.gz
BuildRequires:  golang-packaging
BuildRequires:  coreutils
%if 0%{?rhel}
BuildRequires:  golang >= 1.19
%else
BuildRequires:  golang(API) = 1.19
%endif
BuildRequires:  rsyslog

BuildRequires:       gpgme
BuildRequires:       device-mapper-devel
BuildRequires:       libbtrfs-devel
BuildRequires:       libgpgme-devel


%description
Tools for managing uyuni container.

%prep
%autosetup
tar -zxf %{SOURCE1}


%build
export GOFLAGS=-mod=vendor
%goprep %{provider_prefix}
mkdir -p bin
ADM_PATH=%{provider_prefix}/uyuniadm/shared/utils

tag=%{!?_default_tag:latest}
%if "%{?_default_tag}" != ""
    tag='%{_default_tag}'
%endif

image=registry.opensuse.org/uyuni/server
%if "%{?_default_image}" != ""
  image='%{_default_image}'
%endif

chart=oci://registry.opensuse.org/uyuni/server
%if "%{?_default_chart}" != ""
  chart='%{_default_chart}'
%endif

go build \
    -ldflags "-X ${ADM_PATH}.DefaultImage=%{_image} -X ${ADM_PATH}.DefaultTag=%{_tag} -X ${ADM_PATH}.DefaultChart=%{_chart}" \
    -o ./bin ./...

%install
install -m 0755 -vd %{buildroot}%{_bindir}
install -m 0755 -vp ./bin/* %{buildroot}%{_bindir}/

%gofilelist

%define _release_dir  %{_builddir}/%{project}-%{version}/release

%files

%defattr(-,root,root)
%doc README.md
%license LICENSE

%{_bindir}/uyuniadm
%{_bindir}/uyunictl

%changelog
