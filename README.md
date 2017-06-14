# vSphere Metrics Exporter

[![Build Status](https://travis-ci.org/brandonweeks/vsphere_exporter.svg?branch=master)](https://travis-ci.org/brandonweeks/vsphere_exporter)

[Prometheus](https://prometheus.io/) exporter for VMware vSphere metrics.

## Getting

```
go get github.com/brandonweeks/vsphere_exporter
```

## Building

```
cd $GOPATH/src/github.com/brandonweeks/vsphere_exporter
make
```

## Configuration

### Configuration file

Configuration file in yml format.

```
---
hostname: "vcenter.example.com"
username: "administrator@vsphere.local"
password: ""
datacenter: ""
```

### Environment Variables

Name               | Description
-------------------|------------
VSPHERE_HOSTNAME   | vSphere hostname
VSPHERE_USERNAME   | vSphere username @vsphere.local
VSPHERE_PASSWORD   | vSphere password
VSPHERE_DATACENTER | vSphere datacenter

## Dashboards

TODO
