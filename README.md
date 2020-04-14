# Dnsmasq-controller

A Dnsmasq-controller for Kubernetes, implemented in go using [kubebuilder](https://kubebuilder.io/).

## Status

![GitHub](https://img.shields.io/badge/status-alpha-blue?style=for-the-badge)
![GitHub](https://img.shields.io/github/license/kristofferahl/healthchecksio-operator?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kristofferahl/healthchecksio-operator?style=for-the-badge)

## Supported resources
- DnsmasqOptionSet

## Example
```yaml
---
apiVersion: dnsmasq.kvaps.cf/v1alpha1
kind: DnsmasqOptionSet
metadata:
  name: node1
spec:
  controller: ""
  options:
  - key: host-record
    value: node1,10.9.8.7
  - key: dhcp-host
    value: set:ltsp1,a0:1d:48:b5:ae:a8,a0:1d:48:b5:ae:a9,id:*,10.9.8.7,node1,infinite
```

### Configuration

| Flag                      | Type   | Required | Description                                                                                                                             |
|---------------------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `-cleanup`                | bool   | false    | Cleanup Dnsmasq config directory before start.                                                                                          |
| `-conf-dir`               | string | false    | Dnsmasq config directory for write configuration to. (default "/etc/dnsmasq.d")                                                         |
| `-controller`             | string | false    | Name of the controller this controller satisfies. (default "")                                                                          |
| `-development`            | bool   | false    | Run the controller in development mode.                                                                                                 |
| `-enable-leader-election` | bool   | false    | Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.                   |
| `-kubeconfig`             | string | false    | Paths to a kubeconfig. Only required if out-of-cluster.                                                                                 |
| `-log-level`              | string | false    | The log level used by the operator. (default "info")                                                                                    |
| `-metrics-addr`           | string | false    | The address the metric endpoint binds to. (default ":8080")                                                                             |
| `-sync-delay`             | int    | false    | Time in seconds to syncronise Dnsmasq configuration. (default 1)                                                                        |
| `-watch-namespace`        | string | false    | Namespace the controller watches for updates to Kubernetes objects. All namespaces are watched if this parameter is left empty.         |
| `--`                      | array  | false    | Additional command line arguments for Dnsmasq may be specified after `--` (read [dnsmasq-man] for more details)                         |

[dnsmasq-man]: http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html

## Development

### Pre-requisites
- [Go](https://golang.org/) 1.13 or later
- [Kubebuilder](https://kubebuilder.io/) 2.3.1
- [Kubernetes](https://kubernetes.io/) cluster

### Getting started
```bash
make install
make run
```

### Running tests
```bash
make test
```
