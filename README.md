# vsphere_exporter
A Prometheus exporter for vSphere

## Usage

 ./vsphere_exporter-master --help
Usage of ./vsphere_exporter-master:
-config.file string
vsphere_exporter configuration file name.
-log.level value
Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal, panic].
-web.listen-address string
Address on which to expose metrics and web interface. (default ":9155")
-web.telemetry-path string
Path under which to expose Prometheus metrics. (default "/metrics")

## Modification

The function `ignoreMetric` has an array, that contains metrics exported by vcenter, that will not be exported by this exporter, to reduce the number of time series.
Should be adjusted to own needs. (TODO: Load it via config file)

## Installing

The vsphereExporter.service is an example systemd-service, that restarts the exporter always after it exits (e.g. after a conenction loss to the vcenter).

## Hints

 - Do not call the exporter to often, exporting all metrics take some time and is done live (might need to increase the timeout)
 
