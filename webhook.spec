# Version is unix timestamp of last commit, or
# tag if it exists
%define git_timestamp %(git log --format=format:git%ct)
%define git_tag %(git tag -l --contains HEAD)
%define git_ver %(printf %s "%{git_tag}" "%{git_timestamp}" | sed 1q)

Name: webhook
Version: %{git_ver}
Release: 1
License: MIT
Summary: tcpserver-style Github webhook receiver
Group: System Environment/Daemons
BuildRoot: %{_buildroot}
Source0: webhook.go

%description
webhook listens on a given port for github webhook
payloads, and runs a command for each request.

%prep
rm -rf $RPM_BUILD_DIR/%{name}-%{version}
mkdir -p $RPM_BUILD_DIR/%{name}-%{version}
cp $RPM_SOURCE_DIR/webhook.go $RPM_BUILD_DIR/%{name}-%{version}

%build
rm -rf $RPM_BUILD_ROOT
go build -o webhook

%install
mkdir -p $RPM_BUILD_ROOT%{_bindir}
install -m755 webhook $RPM_BUILD_ROOT%{_bindir}/webhook

%files
%defattr(-,root,root,-)
%{_bindir}/webhook
%{_mandir}/man1/

%changelog
* Fri Nov 15 2013 David Arroyo <darroyo@constantcontact.com> 0.5-1
- Initial build
