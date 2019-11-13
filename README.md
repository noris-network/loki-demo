# [loki-demo](https://github.com/noris-network/loki-demo)

This is a demo of [Grafana Loki][loki], showcasing how to search logs and
export metrics over them.

## Getting started

Clone the repo:

```shell
git clone git@github.com:noris-network/loki-demo.git
cd loki-demo
```

### Run locally

Prerequisites:

- [Docker][docker]
- [Go 1.13+][go]

With this setup Loki, Promtail and the [log_gen][log_gen] script
will run locally, whereas Prometheus and Grafana will run inside
Docker containers on the [host network][docker-net].

First build `log_gen`:

```shell
make build
```

Then download and extract the binaries in to `bin/`:

```shell
make download
```

Next start the services:

```shell
# Each of the below commands should run in its own terminal window
make run/loki
make run/promtail
make run/docker/up
make run/log_gen
```

## Using it

Access Grafana via http://localhost:3000 and add the following
datasources:

Datasource|Name|URL
---|---|---
Prometheus|Prometheus|http://localhost:9090
Loki|Loki|http://localhost:3100
Prometheus|Loki as Prometheus|http://localhost:3100/loki

> The "Loki as Prometheus" datasource is necessary to use aggregation functions like
> `rate` or `count` over LogQL results.

[loki]: https://github.com/grafana/loki
[docker]: https://docs.docker.com/install/
[go]: https://golang.org/doc/install
[docker-net]: https://docs.docker.com/network/#network-driver-summary
[log_gen]: https://github.com/noris-network/loki-demo/blob/master/main.go