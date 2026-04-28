# SPEC file

%global c_vendor    %{_vendor}
%global gh_owner    %{_owner}
%global gh_project  %{_project}

Name:      %{_package}
Version:   %{_version}
Release:   %{_release}%{?dist}
Summary:   Web-Service Random Password Generator

License:   MIT
URL:       https://github.com/%{gh_owner}/%{gh_project}


Provides:  %{gh_project} = %{version}

%description
Web-Service Random Password Generator

%build
#(cd %{_current_directory} && make build)

%install
rm -rf $RPM_BUILD_ROOT
(cd %{_current_directory} && make install DESTDIR=$RPM_BUILD_ROOT)
rm -f $RPM_BUILD_ROOT/etc/passwd

%clean
rm -rf $RPM_BUILD_ROOT

%files
%attr(0755,root,root) %{_binpath}/%{_project}
%attr(0644,root,root) %{_initpath}/%{_project}.service
%attr(-,root,root) %{_docpath}
%attr(0644,root,root) %{_manpath}/%{_project}.1.gz
%docdir %{_docpath}
%docdir %{_manpath}
%config(noreplace) %{_configpath}*

%changelog
* Fri Dec 04 2026 Nicola Asuni <info@tecnick.com> 1.0.0-1
- Initial Commit

