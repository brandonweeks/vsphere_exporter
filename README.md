# vSphere Metrics Exporter

[![Build Status](https://travis-ci.org/brandonweeks/vsphere_exporter.svg?branch=master)](https://travis-ci.org/brandonweeks/vsphere_exporter)

[Prometheus](https://prometheus.io/) exporter for VMware vSphere metrics.

## Getting

```
$ go get github.com/brandonweeks/vsphere_exporter
```

## Building

```
$ cd $GOPATH/src/github.com/brandonweeks/vsphere_exporter
$ make
```

## Docker

```
$ cd $GOPATH/src/github.com/brandonweeks/vsphere_exporter
$ make docker
$ docker run -d -v ${pwd}/vsphere_exporter.yml:/vsphere_exporter.yml \
	--name vsphere_exporter -p 9155:9155 brandonweeks/vsphere_exporter:master \
	/go/bin/app -config.file /vsphere_exporter.yml
```

## Configuration

### Configuration file

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

### Prometheus


Add a block to the `scrape_configs` of your prometheus.yml config file:

```
scrape_configs:
...

- job_name: vsphere_exporter
  scrape_interval: 30s
  scrape_timeout: 25s
  static_configs:
  - targets: ['localhost:9121']

...
```
and adjust the host name accordingly.

## Dashboards

TODO
