FROM golang:1.8-onbuild
MAINTAINER Brandon Weeks <bweeks@google.com>

COPY vsphere_exporter /bin/vsphere_exporter

EXPOSE     9155
ENTRYPOINT [ "/bin/vsphere_exporter" ]
