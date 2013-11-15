Name: webhook
Version: 0.5
Release: 1
License: MIT
Summary: tcpserver-style Github webhook receiver
Group: System Environment/Daemons
BuildRoot: %{_buildroot}

%description
webhook listens on a given port for github webhook
payloads, and runs a command for each request.

%build

rm -rf $RPM_BUILD_ROOT
go build -o webhook

%install

mkdir -p $RPM_BUILD_ROOT%{_bindir}
mkdir -p $RPM_BUILD_ROOT%{_mandir}/man1
install -m755 webhook $RPM_BUILD_ROOT%{_bindir}/webhook
install -m644 webhook.1 $RPM_BUILD_ROOT%{_mandir}/man1/webhook.1

%files
%defattr(-,root,root,-)
%{_bindir}/webhook
%{_mandir}/man1/

%changelog
* Fri Nov 15 2013 David Arroyo <darroyo@constantcontact.com> 0.5-1
- Initial build
