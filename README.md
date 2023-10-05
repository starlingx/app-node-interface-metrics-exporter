# app-node-interface-metrics-exporter
app-node-interface-metrics-exporter flux cd app

####  Top Level Directory Structure
```bash
├── app-node-interface-metrics-exporter    # Root Folder
│   ├── debian_build_layer.cfg
│   ├── debian_iso_image.inc
│   ├── debian_pkg_dirs
│   ├── metrics-exporter-api              # Go code  for  api server which will expose  Metrics for  NIC .
│   ├── python3-k8sapp-node-interface-metrics-exporter    # lifecycle managemnt code  to support  flux apps
│   ├── README.md
│   ├── requirements.txt
│   ├── stx-node-interface-metrics-exporter-helm      # helm Package manager  for the app
│   ├── test-requirements.txt
│   └── tox.ini
```

> all command related to go code should run from `app-node-interface-metrics-exporter/metrics-exporter-api/docker/metrics-exporter-api`
### About metrics-exporter-api
It is Simple Http Server which reads network interface (PCI devices) and exposes following API's

1) http://{hostname}:{port}/metrics -- all node metrics in [OpenMetrics] format
2) http://{hostname}:{port}/metrics/device/{DeviceName} -- specified interface metrics identified by device name in [OpenMetrics] format
3) http://{hostname}:{port}/metrics/pci-addr/{PciAddr} -- specified interface metrics identified by PCI address in [OpenMetrics] format
4) http://{hostname}:{port}/json/metrics -- all node metrics in [JSON] format
5) http://{hostname}:{port}/json/metrics/device/{DeviceName} -- specified interface metrics identified by device name in [JSON] format
6) http://{hostname}:{port}/json/metrics/pci-addr/{PciAddr} -- specified interface metrics identified by PCI address in [JSON] format

#### Makefile Support
```bash
$ make
help:                Show this help.
install_dep:         install go dependency
run:                 run app on host machine
test:                run go unit test
testcov:             run go coverage test
vet:                 run go vet
lint:                run go lint
build_linux:         Build application
build_image:         Build docker image
```
> `$ make run ` will start go dev http server

#### Run time Options / params for the metrics-exporter-app usage
| Options | Help |
| ------ | ------ |
| `-log.file`           | Log file name (default "node_metrics_api.log") |
| `-log.level`          | log level. (default "info"). Valid options trace, debug, info, warning, error, fatal and panic |
| `-log.file`           | Log file name (default "node_metrics_api.log") |
| `-web.listen-address` | Port to listen on for web interface. (default ":9110") |
| `-path.sys`           | mounted path for host /sys inside container (default "/sys") |

#### Local / Devlopment Set UP for metrics-exporter-api
`pre requisite go 1.21.0`

#### Container image reference for helm
[Dockerfile](/metrics-exporter-api/debian/Dockerfile)

[Build Reference](/metrics-exporter-api/debian/metrics-exporter-api.stable_docker_image)

#### References
[StarlingX](https://www.starlingx.io/)

[How to add a FluxCD App to Starlingx](https://wiki.openstack.org/wiki/StarlingX/Containers/HowToAddNewFluxCDAppInSTX)

[OpenMetrics](https://openmetrics.io/)

[JSON]: <https://www.json.org/>
[OpenMetrics]: <https://openmetrics.io/>