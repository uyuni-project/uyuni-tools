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
%global org             mbussolotto
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
%if 0%{?rhel}
BuildRequires:  golang >= 1.18
%else
BuildRequires:  golang(API) = 1.18
%endif
BuildRequires:  rsyslog

Requires:       gpgme
Requires:       libbtrfs-devel
Requires:       libassuan
Requires:       libgpgme-devel


%description
Tools for managing uyuni container.

%prep
%autosetup
tar -zxf %{SOURCE1}


%build
export GOFLAGS=-mod=vendor
%goprep %{provider_prefix}
%gobuild ...

%install
%goinstall
%gosrc

%gofilelist

%define _release_dir  %{_builddir}/%{project}-%{version}/release

%files

%defattr(-,root,root)
%doc README.md
%license LICENSE
%{_bindir}/uyuni-tools

%config(noreplace) %{_sysconfdir}/uyuni-tools/options.json

%changelog
