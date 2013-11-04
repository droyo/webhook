Name: webhook
Version: 1.0
Release: 1
License: MIT
Group: Applications/System
URL: https://github.roving.com/darroyo/webhook
Summary: Webhook listener in the spirit of inetd
Source0: webhook.go
Source1: webhook.1
BuildRoot: %{_buildroot}

%description
webhook executes a program whenever it receives a webhook request
from Github.  The child program is run with repository information
set in its environment, and with the json commit data written to
its standard input. webhook guarantees that only one child program
is run at a time.

%build
rm -rf $RPM_BUILD_ROOT
go build -o webhook

%install
mkdir -p $RPM_BUILD_ROOT%{_bindir}
mkdir -p $RPM_BUILD_ROOT%{_mandir}/man1
install -m755 webhook $RPM_BUILD_ROOT%{_bindir}/webhook
install -m644 webhook.1 $RPM_BUILD_ROOT%{_mandir}/man1/webhook

%files
%defattr(-,root,root,-)
%{_bindir}/webhook
%{_mandir}/man1/

%changelog
* Tue Nov 04 2013 David Arroyo <darroyo@constantcontact.com>
- Initial build
