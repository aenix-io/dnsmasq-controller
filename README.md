# Dnsmasq-controller

A Dnsmasq-controller for Kubernetes, implemented in go using [kubebuilder](https://kubebuilder.io/).

## Status

![GitHub](https://img.shields.io/badge/status-alpha-blue?style=for-the-badge)
![GitHub](https://img.shields.io/github/license/kristofferahl/healthchecksio-operator?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kristofferahl/healthchecksio-operator?style=for-the-badge)

## Supported resources

- DnsmasqOptionSet
- DnsmasqHostSet
- DnsmasqDhcpHostSet
- DnsmasqDhcpOptionSet

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
  - key: txt-record
    value: example.com,"v=spf1 a -all"
---
apiVersion: dnsmasq.kvaps.cf/v1alpha1
kind: DnsmasqHostSet
metadata:
  name: node1
spec:
  controller: ""
  hosts:
  - ip: "10.9.8.7"
    hostnames:
    - "foo.local"
    - "bar.local"
  - ip: "10.1.2.3"
    hostnames:
    - "foo.remote"
    - "bar.remote"
---
apiVersion: dnsmasq.kvaps.cf/v1alpha1
kind: DnsmasqDhcpHostSet
metadata:
  name: node1
spec:
  controller: ""
  hosts:
  - ip: 10.9.8.7
    macs:
    - 94:57:a5:d3:b6:f2
    - 94:57:a5:d3:b6:f3
    clientIDs:
    - "*"
    setTags:
    - hp
    tags:
    - ltsp1
    hostname: node1
    leaseTime: infinite
    ignore: false
---
apiVersion: dnsmasq.kvaps.cf/v1alpha1
kind: DnsmasqDhcpOptionSet
metadata:
  name: ltsp1
spec:
  controller: ""
  options:
  - type: option
    key: ntp-server
    values:
    - 192.168.0.4
    tags:
    - ltsp1
```

	flag.BoolVar(&config.EnableDNS, "dns", false, "Enable DNS Service and DNS configuration discovery.")
	flag.BoolVar(&config.EnableDHCP, "dhcp", false, "Enable DHCP Service and DHCP configuration discovery.")

### Configuration

| Flag                      | Type   | Required | Description                                                                                                                             |
|---------------------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `-cleanup`                | bool   | false    | Cleanup Dnsmasq config directory before start.                                                                                          |
| `-conf-dir`               | string | false    | Dnsmasq config directory for write configuration to. (default "/etc/dnsmasq.d")                                                         |
| `-controller`             | string | false    | Name of the controller this controller satisfies. (default "")                                                                          |
| `-development`            | bool   | false    | Run the controller in development mode.                                                                                                 |
| `-dhcp`                   | bool   | false    | Enable DHCP Service and configuration discovery.                                                                                        |
| `-dns`                    | bool   | false    | Enable DNS Service and configuration discovery.                                                                                         |
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
