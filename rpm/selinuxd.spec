Name:          selinuxd
Version:       0.1

%global goipath github.com/JAORMX/selinuxd
%global commit  a748373d23f48245f8db9672fd923c2bdd1859ec
%gometa

Release:       1%{?dist}
Summary:       TBD

License:       ASL 2.0
URL:           %{gourl}
Source0:       selinuxd.tar.gz

BuildRequires: go-rpm-macros
BuildRequires: make
BuildRequires: libsemanage-devel
# not strictly a BR but we copy the base policies from here
BuildRequires: udica
# to be able to use the systemd scriptlets
BuildRequires: systemd

# this is probably a bug in libsemanage; it forks and execs
# binaries from policycoreutils..eek
Requires: policycoreutils
# to install base policy elements
Requires: container-selinux

%description
TBD

%prep
%setup -q

%build
%{__make}

%install
%{__make} PREFIX=%{buildroot}%{_prefix} ETCDIR=%{buildroot}%{_sysconfdir} LOCALSTATEDIR=%{buildroot}%{_localstatedir}/run \
    install \
    install.systemd

# This is hack until we split off the base policies from udica
# which requires python to a standalone package
%{__install} -Z -m 644 /usr/share/udica/templates/* %{buildroot}%{_sysconfdir}/selinux.d/

%files
%{_bindir}/selinuxdctl
%{_unitdir}/%{name}.socket
%{_unitdir}/%{name}.service
%dir %attr(700, root, root) %{_sysconfdir}/selinux.d/
%{_sysconfdir}/selinux.d/*.cil

%post
%systemd_post selinuxd.service

%preun
%systemd_preun selinuxd.service

%postun
%systemd_postun_with_restart selinuxd.service
