# Dnsmasq-controller

A Dnsmasq-controller for Kubernetes, implemented in go using [kubebuilder](https://kubebuilder.io/).

## Status

![GitHub](https://img.shields.io/badge/status-beta-blue?style=for-the-badge)
![GitHub](https://img.shields.io/github/license/kristofferahl/healthchecksio-operator?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kristofferahl/healthchecksio-operator?style=for-the-badge)

## Supported resources

- DnsmasqOptions
- DnsHosts
- DhcpHosts
- DhcpOptions


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

## Installation

```bash
# CRDs
kubectl apply -k config/crd/bases

# RBAC
kubectl apply -k config/rbac

# DNS-server (for infra.example.org)
kubectl apply -k config/dns-server

# DHCP-server
kubectl apply -k config/dhcp-server

# Add dnsmasq role to your nodes
kubectl label node <node1> <node2> <node3> node-role.kubernetes.io/dnsmasq=
```

## Examples

Global DHCP-configuration:

```yaml
---
apiVersion: dnsmasq.kvaps.cf/v1beta1
kind: DhcpOptions
metadata:
  name: default-network-configuration
spec:
  controller: ""
  options:
  - key: option:router
    values: [192.168.67.1]
  - key: option:dns-server
    values: [192.168.67.1]
  - key: option:domain-name
    values: [infra.example.org]
  - key: option:domain-search
    values: [infra.example.org]
---
apiVersion: dnsmasq.kvaps.cf/v1beta1
kind: DnsmasqOptions
metadata:
  name: default-matchers
spec:
  controller: ""
  options:
  - key: dhcp-range
    values: [192.168.67.0,static,infinite]
  - key: dhcp-match
    values: [set:iPXE,"175","39"]
  - key: dhcp-match
    values: [set:X86PC,option:client-arch,"0"]
  - key: dhcp-match
    values: [set:X86-64_EFI,option:client-arch,"7"]
  - key: dhcp-match
    values: [set:X86-64_EFI,option:client-arch,"9"]
```

Global DNS-configuration:

```yaml
---
apiVersion: dnsmasq.kvaps.cf/v1beta1
kind: DnsmasqOptions
metadata:
  name: global-dns
spec:
  controller: ""
  options:
  - key: srv-host
    values: [_kerberos-master._tcp.infra.example.org,freeipa.example.org,"88"]
  - key: srv-host
    values: [_kerberos-master._udp.infra.example.org,freeipa.example.org,"88"]
  - key: srv-host
    values: [_kerberos._tcp.infra.example.org,freeipa.example.org,"88"]
  - key: srv-host
    values: [_kerberos._udp.infra.example.org,freeipa.example.org,"88"]
  - key: srv-host
    values: [_kpasswd._tcp.infra.example.org,freeipa.example.org,"464"]
  - key: srv-host
    values: [_kpasswd._udp.infra.example.org,freeipa.example.org,"464"]
  - key: srv-host
    values: [_ldap._tcp.infra.example.org,freeipa.example.org,"389"]
  - key: srv-host
    values: [_ntp._udp.infra.example.org,129.6.15.28,"123"]
  - key: srv-host
    values: [_ntp._udp.infra.example.org,129.6.15.29,"123"]
  - key: txt-record
    values: [_kerberos.infra.example.org,EXAMPLE.ORG]
```

Netboot-server configuration with tag `ltsp1`:

```yaml
---
apiVersion: dnsmasq.kvaps.cf/v1beta1
kind: DhcpOptions
metadata:
  name: ltsp1
spec:
  controller: ""
  options:
  - key: option:server-ip-address
    tags: [ltsp1]
    values: [192.168.67.11]
  - key: option:tftp-server
    tags: [ltsp1]
    values: [ltsp1]
  - key: option:bootfile-name
    tags: [ltsp1,X86PC]
    values: [ltsp/grub/i386-pc/core.0]
  - key: option:bootfile-name
    tags: [ltsp1,X86-64_EFI]
    values: [ltsp/grub/x86_64-efi/core.efi]
```

DHCP-client for network booting using assigned tag `ltsp1`:

```yaml
---
apiVersion: dnsmasq.kvaps.cf/v1beta1
kind: DhcpHosts
metadata:
  name: netboot-client
spec:
  controller: ""
  hosts:
  - ip: 192.168.67.20
    macs:
    - 94:57:a5:d3:b6:f2
    - 94:57:a5:d3:b6:f3
    clientIDs: ["*"]
    setTags: [ltsp1]
    hostname: node1
    leaseTime: infinite
```

Add A, AAAA and PTR records to the DNS:

```yaml
---
apiVersion: dnsmasq.kvaps.cf/v1beta1
kind: DnsHosts
metadata:
  name: netboot-client
spec:
  controller: ""
  hosts:
  - ip: 192.168.67.20
    hostnames:
    - node1
    - node1.infra.example.org
```

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
