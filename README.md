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

Run `make install` to build `ts_gen` and download Loki and
Promtail:

```shell
make install
```

Next start the services:

> Each of the below commands should be run in a separate
> terminal window.

```shell
make run/loki
make run/promtail
make run/docker/up
make run/log_gen
```

If you have an AWS S3 bucket, you can store loki's chunks in it with.

```shell
make run/loki/s3 ACCESSKEY=<your aws access key>  SECRETKEY=<aws secret key> \
     S3ENDPOINT=<s3 endoint>  BUCKETNAME=<bucket name>
```

## Using it

Access Grafana via http://localhost:3000 and add the following
datasources:

Datasource|Name|URL
---|---|---
Prometheus|Prometheus|http://localhost:9090
Loki|Loki|http://localhost:3100
Prometheus|Loki as Prometheus|http://localhost:3100/loki

> As of Grafana 6.4 the `Loki as Prometheus` datasource is necessary to use aggregation functions like
> `rate` or `count` over LogQL results.

![Datasources][ds_pic]

You can then import the sample dashboard in `dashboards/log_gen.json`:

![log_gen Dashboards][dashboard_example]

Going on **Explore** on the left hand side lets you evaluate
the logs via [LogQL][logql]:

![Example LogQL query][query1]

```
{job="demo_log", service="api", level="error"} |= "cpu"
```

> As of Grafana 6.4, LogQL functions need to be sent agaist the
> `Loki as Prometheus` datasource.

![Example LogQL query with functions][query2]

```
sum by (handler) (rate({job="demo_log", handler!=""})[5m])
```

## Cleanup

First stop all services that run in the foreground:

```shell
CTRL+c
```

Then shutdown the docker-compose stack:

```shell
make run/docker/down
```

Running `make clean` will remove the log file created by `log_gen`, all
binaries, created Docker volumes and all data created by Loki and Promtail:

```shell
make clean
```

[loki]: https://github.com/grafana/loki
[docker]: https://docs.docker.com/install/
[go]: https://golang.org/doc/install
[docker-net]: https://docs.docker.com/network/#network-driver-summary
[log_gen]: https://github.com/noris-network/loki-demo/blob/master/main.go
[ds_pic]: static/ds.png
[dashboard_example]: static/dashboard_example.png
[logql]: https://github.com/grafana/loki/blob/master/docs/logql.md
[query1]: static/log_query_1.png
[query2]: static/log_query_2.png

