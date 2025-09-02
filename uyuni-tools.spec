#
# spec file for package uyuni-tools
#
# Copyright (c) 2024 SUSE LLC
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
%{!?productprettyname: %global productprettyname Uyuni}

%global namespace       registry.opensuse.org/uyuni

%if "%{productprettyname}" == "Uyuni"
%if 0%{?suse_version} >= 1600 || 0%{?sle_version} >= 150400 || 0%{?rhel} >= 8 || 0%{?fedora} >= 37 || 0%{?debian} >= 12 || 0%{?ubuntu} >= 2004
%define adm_build    1
%else
%define adm_build    0
%endif
%else
%if 0%{?suse_version} >= 1600 || 0%{?sle_version} >= 150400
%define adm_build 1
%else
%define adm_build 0
%endif
%endif

%if 0%{?debian}
# Don't build kubernetes support for Debian since go is too old (<1.21) there.
%define _uyuni_tools_tags nok8s
%endif

%define name_adm mgradm
%define name_ctl mgrctl
%define name_pxy mgrpxy

# Completion files
%if 0%{?debian} || 0%{?ubuntu}
%define _zshdir %{_datarootdir}/zsh/vendor-completions
%else
%define _zshdir %{_datarootdir}/zsh/site-functions
%endif
# 0%{?debian} || 0%{?ubuntu}

Name:           %{project}
Version:        5.1.20
Release:        0
Summary:        Tools for managing %{productprettyname} container
License:        Apache-2.0
Group:          System/Management
URL:            https://%{provider_prefix}
Source0:        %{name}-%{version}.tar.gz
Source1:        vendor.tar.gz
BuildRequires:  bash-completion
BuildRequires:  coreutils
%if 0%{?debian} || 0%{?ubuntu}
BuildRequires:  gettext
%endif
# 0%{?debian} || 0%{?ubuntu}

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
BuildRequires:  fish
%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}

BuildRequires:  zsh
# Get the proper Go version on different distros
%if 0%{?suse_version}
BuildRequires:  golang(API) >= 1.22
%endif
# 0%{?suse_version}

%if 0%{?ubuntu}
%define go_version      1.22
BuildRequires:  golang-%{go_version}
%endif
# 0%{?ubuntu}

%if 0%{?debian}
BuildRequires:  golang >= 1.19
%endif
# 0%{?debian}

%if 0%{?fedora} || 0%{?rhel}
BuildRequires:  golang >= 1.21
%endif
# 0%{?fedora} || 0%{?rhel}

%description
Tools for managing %{productprettyname} container.

%if %{adm_build}
%package -n %{name_adm}
Summary:        Command line tool to install and update %{productprettyname}
%if 0%{?suse_version}
Requires:       (aardvark-dns if podman)
Requires:       (netavark if podman)
%endif
# 0%{?suse_version}
%if "%{_vendor}" != "debbuild"
Requires:       (podman >= 4.5.0 if podman)
%endif

%description -n %{name_adm}
%{name_adm} is a convenient tool to install and update %{productprettyname} components as containers running
either on Podman or a Kubernetes cluster.

%package -n %{name_pxy}
Summary:        Command line tool to install and update %{productprettyname} proxy
Obsoletes:      uyuni-proxy-systemd-services
%if 0%{?suse_version}
Requires:       (aardvark-dns if podman)
Requires:       (netavark if podman)
%endif
# 0%{?suse_version}

%description -n %{name_pxy}
%{name_pxy} is a convenient tool to install and update %{productprettyname} proxy components as containers
running either on Podman or a Kubernetes cluster.

%package -n %{name_adm}-bash-completion
Summary:        Bash Completion for %{name_adm}
Group:          System/Shells
Requires:       %{name_adm} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_adm} and bash-completion)
%else
Supplements:    bash-completion
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_adm}-bash-completion
Bash command line completion support for %{name_adm}.

%package -n %{name_adm}-zsh-completion
Summary:        Zsh Completion for %{name_adm}
Group:          System/Shells
Requires:       %{name_adm} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_adm} and zsh)
%else
Supplements:    zsh
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_adm}-zsh-completion
Zsh command line completion support for %{name_adm}.

%package -n %{name_pxy}-bash-completion
Summary:        Bash Completion for %{name_pxy}
Group:          System/Shells
Requires:       %{name_pxy} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_pxy} and bash-completion)
%else
Supplements:    bash-completion
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_pxy}-bash-completion
Bash command line completion support for %{name_pxy}.

%package -n %{name_pxy}-zsh-completion
Summary:        Zsh Completion for %{name_pxy}
Group:          System/Shells
Requires:       %{name_pxy} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_pxy} and zsh)
%else
Supplements:    zsh
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_pxy}-zsh-completion
Zsh command line completion support for %{name_pxy}.

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
%package -n %{name_adm}-fish-completion
Summary:        Fish Completion for %{name_adm}
Group:          System/Shells
Requires:       %{name_adm} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_adm} and fish)
%else
Supplements:    fish
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_adm}-fish-completion
Fish command line completion support for %{name_adm}.

%package -n %{name_pxy}-fish-completion

Summary:        Fish Completion for %{name_pxy}
Group:          System/Shells
Requires:       %{name_pxy} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_pxy} and fish)
%else
Supplements:    fish
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_pxy}-fish-completion
Fish command line completion support for %{name_pxy}.

%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}

%endif
# %{adm_build}

%package -n %{name_ctl}
Summary:        Command line tool to perform day-to-day operations on %{productprettyname}

%description -n %{name_ctl}
%{name_ctl} is a tool helping with daily tasks on %{productprettyname} components running as containers
either on Podman or a Kubernetes cluster.

%package -n %{name_ctl}-bash-completion
Summary:        Bash Completion for %{name_ctl}
Group:          System/Shells
Requires:       %{name_ctl} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_ctl} and bash-completion)
%else
Supplements:    bash-completion
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_ctl}-bash-completion
Bash command line completion support for %{name_ctl}.

%package -n %{name_ctl}-zsh-completion
Summary:        Zsh Completion for %{name_ctl}
Group:          System/Shells
Requires:       %{name_ctl} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_ctl} and zsh)
%else
Supplements:    zsh
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_ctl}-zsh-completion
Zsh command line completion support for %{name_ctl}.

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
%package -n %{name_ctl}-fish-completion
Summary:        Fish Completion for %{name_ctl}
Group:          System/Shells
Requires:       %{name_ctl} = %{version}
BuildArch:      noarch
%if 0%{?suse_version} >= 150000
Supplements:    (%{name_ctl} and fish)
%else
Supplements:    fish
%endif
# 0%{?suse_version} >= 150000

%description -n %{name_ctl}-fish-completion
Fish command line completion support for %{name_ctl}.
%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}

# Only SUSE distros have a -lang packages, for the others they
# will all be in the correspdonding tool package.
%if 0%{?suse_version} || 0%{?sle_version}
%{lang_package -n %{name_ctl}}
%{lang_package -n %{name_pxy}}

%if %{adm_build}
%{lang_package -n %{name_adm}}
%endif
# %{adm_build}

%endif
# 0%{?suse_version} || 0%{?sle_version}

%prep
%autosetup
tar -zxf %{SOURCE1}

%build
%ifarch i386
%if 0%{?debian}
# Disable CGO build for debian 32 bits to avoid cross-compilation
export CGO_ENABLED=0
%endif
%endif

export GOFLAGS=-mod=vendor
mkdir -p bin
UTILS_PATH="%{provider_prefix}/shared/utils"

tag=%{!?_default_tag:latest}
%if "%{?_default_tag}" != ""
    tag='%{_default_tag}'
%endif
# "%{?_default_tag}" != ""

pull_policy=%{!?_default_pull_policy:Always}
%if "%{?_default_pull_policy}" != ""
    pull_policy='%{_default_pull_policy}'
%endif
# "%{?_default_pull_policy}" != ""

namespace=%{namespace}
helm_registry=%{namespace}
%if "%{?_default_namespace}" != ""
  # Set both container and helm chart namespaces as this can be the same value
  namespace='%{_default_namespace}'
  helm_registry='%{_default_namespace}'
%endif
# "%{?_default_namespace}" != ""

# We may have additional config for helm registry as the path is different in OBS devel projects
%if "%{?_default_helm_registry}" != ""
  helm_registry='%{_default_helm_registry}'
%endif
# "%{?_default_helm_registry}" != ""

go_tags=""
%if "%{?_uyuni_tools_tags}" != ""
  go_tags="-tags %{_uyuni_tools_tags}"
%endif
# "%{?_uyuni_tools_tags}" != ""

go_path=""
%if 0%{?ubuntu}
  go_path=/usr/lib/go-%{go_version}/bin/
%else
  %if "%{?_go_bin}" != ""
    go_path='%{_go_bin}/'
  %endif
# "%{?_go_bin}" != ""

%endif
# 0%{?ubuntu}

GOLD_FLAGS="-X '${UTILS_PATH}.Version=%{version} for ${tag} image (%{version_details}) (compilation tag: %{_uyuni_tools_tags})' -X ${UTILS_PATH}.LocaleRoot=%{_datadir}/locale"
if test -n "${namespace}"; then
    GOLD_FLAGS="${GOLD_FLAGS} -X ${UTILS_PATH}.DefaultRegistry=${namespace}"
fi

if test -n "${helm_registry}"; then
    GOLD_FLAGS="${GOLD_FLAGS} -X ${UTILS_PATH}.DefaultHelmRegistry=${helm_registry}"
fi

if test -n "${tag}"; then
    GOLD_FLAGS="${GOLD_FLAGS} -X ${UTILS_PATH}.DefaultTag=${tag}"
fi

if test -n "${pull_policy}"; then
    GOLD_FLAGS="${GOLD_FLAGS} -X ${UTILS_PATH}.DefaultPullPolicy=${pull_policy}"
fi

# Workaround for rpm on Fedora and EL clones not able to handle go's compressed debug symbols
# Found compressed .debug_aranges section, not attempting dwz compression
%if 0%{?rhel} >= 8 || 0%{?fedora} >= 38
GOLD_FLAGS="-compressdwarf=false ${GOLD_FLAGS}"
%endif
# 0%{?rhel} >= 8 || 0%{?fedora} >= 38

# Workaround for missing build-id on Fedora
# error: Missing build-id in [...]
%if 0%{?fedora} >= 38
GOLD_FLAGS="-B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') ${GOLD_FLAGS}"
%endif
# 0%{?fedora} >= 38

${go_path}go build ${go_tags} -ldflags "${GOLD_FLAGS}" -o ./bin ./...

%if ! %{adm_build}
rm ./bin/%{name_adm}
rm ./bin/%{name_pxy}
%endif
# ! %{adm_build}

%install
install -m 0755 -vd %{buildroot}%{_bindir}
install -m 0755 -vp ./bin/* %{buildroot}%{_bindir}/

# Generate the machine object files for localizations
./locale/build.sh %{buildroot}%{_datadir}/locale/

%find_lang %{name_ctl}
%if %{adm_build}
%find_lang %{name_adm}
%find_lang %{name_pxy}
%else
rm %{buildroot}%{_datadir}/locale/*/LC_MESSAGES/%{name_adm}.mo
rm %{buildroot}%{_datadir}/locale/*/LC_MESSAGES/%{name_pxy}.mo
%endif
# %{adm_build}

# Completion files
mkdir -p %{buildroot}%{_datarootdir}/bash-completion/completions/
mkdir -p %{buildroot}%{_zshdir}

%{buildroot}/%{_bindir}/%{name_ctl} completion bash > %{buildroot}%{_datarootdir}/bash-completion/completions/%{name_ctl}
%{buildroot}/%{_bindir}/%{name_ctl} completion zsh > %{buildroot}%{_zshdir}/_%{name_ctl}

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
mkdir -p %{buildroot}%{_datarootdir}/fish/vendor_completions.d/
%{buildroot}/%{_bindir}/%{name_ctl} completion fish > %{buildroot}%{_datarootdir}/fish/vendor_completions.d/%{name_ctl}.fish
%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}

%if %{adm_build}

%{buildroot}/%{_bindir}/%{name_adm} completion bash > %{buildroot}%{_datarootdir}/bash-completion/completions/%{name_adm}
%{buildroot}/%{_bindir}/%{name_adm} completion zsh > %{buildroot}%{_zshdir}/_%{name_adm}

%{buildroot}/%{_bindir}/%{name_pxy} completion bash > %{buildroot}%{_datarootdir}/bash-completion/completions/%{name_pxy}
%{buildroot}/%{_bindir}/%{name_pxy} completion zsh > %{buildroot}%{_zshdir}/_%{name_pxy}

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
%{buildroot}/%{_bindir}/%{name_adm} completion fish > %{buildroot}%{_datarootdir}/fish/vendor_completions.d/%{name_adm}.fish
%{buildroot}/%{_bindir}/%{name_pxy} completion fish > %{buildroot}%{_datarootdir}/fish/vendor_completions.d/%{name_pxy}.fish
%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}

%endif
# %{adm_build}

%if %{adm_build}

# mgradm packages files

# Only SUSE distros have a -lang package
%if 0%{?suse_version} || 0%{?sle_version}
%files -n %{name_adm}-lang -f %{name_adm}.lang

%files -n %{name_adm}
%else
%files -n %{name_adm} -f %{name_adm}.lang
%endif
# 0%{?suse_version} || 0%{?sle_version}

%defattr(-,root,root)
%doc README.md
%license LICENSE
%{_bindir}/%{name_adm}

%files -n %{name_adm}-bash-completion
%{_datarootdir}/bash-completion/completions/%{name_adm}

%files -n %{name_adm}-zsh-completion
%{_zshdir}/_%{name_adm}

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
%files -n %{name_adm}-fish-completion
%{_datarootdir}/fish/vendor_completions.d/%{name_adm}.fish
%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}


# mgrpxy packages files

# Only SUSE distros have a -lang package
%if 0%{?suse_version} || 0%{?sle_version}
%files -n %{name_pxy}-lang -f %{name_pxy}.lang

%files -n %{name_pxy}
%else
%files -n %{name_pxy} -f %{name_pxy}.lang
%endif
# 0%{?suse_version} || 0%{?sle_version}

%defattr(-,root,root)
%doc README.md
%license LICENSE
%{_bindir}/%{name_pxy}

%files -n %{name_pxy}-bash-completion
%{_datarootdir}/bash-completion/completions/%{name_pxy}

%files -n %{name_pxy}-zsh-completion
%{_zshdir}/_%{name_pxy}

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
%files -n %{name_pxy}-fish-completion
%{_datarootdir}/fish/vendor_completions.d/%{name_pxy}.fish
%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}

%endif
# %{adm_build}

# mgrctl packages files

# Only SUSE distros have a -lang package
%if 0%{?suse_version} || 0%{?sle_version}
%files -n %{name_ctl}-lang -f %{name_ctl}.lang

%files -n %{name_ctl}
%else
%files -n %{name_ctl} -f %{name_ctl}.lang
%endif
# 0%{?suse_version} || 0%{?sle_version}

%defattr(-,root,root)
%doc README.md
%license LICENSE
%{_bindir}/%{name_ctl}

%files -n %{name_ctl}-bash-completion
%{_datarootdir}/bash-completion/completions/%{name_ctl}

%files -n %{name_ctl}-zsh-completion
%{_zshdir}/_%{name_ctl}

%if 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}
%files -n %{name_ctl}-fish-completion
%{_datarootdir}/fish/vendor_completions.d/%{name_ctl}.fish
%endif
# 0%{?is_opensuse} || 0%{?fedora} || 0%{?debian} || 0%{?ubuntu}

%changelog
