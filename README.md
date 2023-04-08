# Disk Usage Prometheus Exporter

[![Build Status](https://travis-ci.com/dundee/disk_usage_exporter.svg?branch=master)](https://travis-ci.com/dundee/disk_usage_exporter)
[![codecov](https://codecov.io/gh/dundee/disk_usage_exporter/branch/master/graph/badge.svg)](https://codecov.io/gh/dundee/disk_usage_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/dundee/disk_usage_exporter)](https://goreportcard.com/report/github.com/dundee/disk_usage_exporter)
[![Maintainability](https://api.codeclimate.com/v1/badges/74d685f0c638e6109ab3/maintainability)](https://codeclimate.com/github/dundee/disk_usage_exporter/maintainability)
[![CodeScene Code Health](https://codescene.io/projects/14689/status-badges/code-health)](https://codescene.io/projects/14689)

Provides detailed info about disk usage of the selected filesystem path.

Uses [gdu](https://github.com/dundee/gdu) under the hood for the disk usage analysis.

## Demo Grafana dashboard

https://grafana.milde.cz/d/0TfJhs_Mz/disk-usage (credentials: grafana / gdu)

## Usage

```
Usage:
  disk_usage_exporter [flags]

Flags:
  -p, --analyzed-path string   Path where to analyze disk usage (default "/")
  -b, --bind-address string    Address to bind to (default "0.0.0.0:9995")
  -c, --config string          config file (default is $HOME/.disk_usage_exporter.yaml)
  -l, --dir-level int          Directory nesting level to show (0 = only selected dir) (default 2)
  -L, --follow-symlinks        Follow symlinks for files, i.e. show the size of the file to which symlink points to (symlinks to directories are not followed)
  -h, --help                   help for disk_usage_exporter
  -i, --ignore-dirs strings    Absolute paths to ignore (separated by comma) (default [/proc,/dev,/sys,/run,/var/cache/rsnapshot])
  -m, --mode string            Exposition method - either 'file' or 'http' (default "http")
  -f, --output-file string     Target file to store metrics in (default "./disk-usage-exporter.prom")
```

## Example output

```
# HELP node_disk_usage_bytes Disk usage of the directory/file
# TYPE node_disk_usage_bytes gauge
node_disk_usage_bytes{path="/var/cache"} 2.1766144e+09
node_disk_usage_bytes{path="/var/db"} 20480
node_disk_usage_bytes{path="/var/dpkg"} 8192
node_disk_usage_bytes{path="/var/empty"} 4096
node_disk_usage_bytes{path="/var/games"} 4096
node_disk_usage_bytes{path="/var/lib"} 7.554709504e+09
node_disk_usage_bytes{path="/var/local"} 4096
node_disk_usage_bytes{path="/var/lock"} 0
node_disk_usage_bytes{path="/var/log"} 4.247068672e+09
node_disk_usage_bytes{path="/var/mail"} 0
node_disk_usage_bytes{path="/var/named"} 4096
node_disk_usage_bytes{path="/var/opt"} 4096
node_disk_usage_bytes{path="/var/run"} 0
node_disk_usage_bytes{path="/var/snap"} 1.11694848e+10
node_disk_usage_bytes{path="/var/spool"} 16384
node_disk_usage_bytes{path="/var/tmp"} 475136
# HELP node_disk_usage_level_1_bytes Disk usage of the directory/file level 1
# TYPE node_disk_usage_level_1_bytes gauge
node_disk_usage_level_1_bytes{path="/bin"} 0
node_disk_usage_level_1_bytes{path="/boot"} 1.29736704e+08
node_disk_usage_level_1_bytes{path="/etc"} 1.3090816e+07
node_disk_usage_level_1_bytes{path="/home"} 8.7081373696e+10
node_disk_usage_level_1_bytes{path="/lib"} 0
node_disk_usage_level_1_bytes{path="/lib64"} 0
node_disk_usage_level_1_bytes{path="/lost+found"} 4096
node_disk_usage_level_1_bytes{path="/mnt"} 4096
node_disk_usage_level_1_bytes{path="/opt"} 2.979229696e+09
node_disk_usage_level_1_bytes{path="/root"} 4096
node_disk_usage_level_1_bytes{path="/sbin"} 0
node_disk_usage_level_1_bytes{path="/snap"} 0
node_disk_usage_level_1_bytes{path="/srv"} 4.988928e+06
node_disk_usage_level_1_bytes{path="/tmp"} 1.3713408e+07
node_disk_usage_level_1_bytes{path="/usr"} 1.8109427712e+10
node_disk_usage_level_1_bytes{path="/var"} 2.5156793856e+10
```

## Example Prometheus queries

Disk usage of `/var` directory:

```
sum(node_disk_usage_bytes{path=~"/var.*"})
```

## Example config files

`~/.disk_usage_exporter.yaml`:
```yaml
analyzed-path: /
bind-address: 0.0.0.0:9995
dir-level: 2
ignore-dirs:
- /proc
- /dev
- /sys
- /run
```

`~/.disk_usage_exporter.yaml`:
```yaml
analyzed-path: /
mode: file
output-file: ./disk-usage-exporter.prom
dir-level: 2
ignore-dirs:
- /proc
- /dev
- /sys
- /run
```

## Prometheus scrape config

Disk usage analysis can be resource heavy.
Set the `scrape_interval` and `scrape_timeout` according to the size of analyzed path.

```yaml
scrape_configs:
  - job_name: 'disk-usage'
    scrape_interval: 5m
    scrape_timeout: 20s
    static_configs:
    - targets: ['localhost:9995']
```

## Dump to file

The official `node-exporter` allows to specify a folder which contains additional metric files through a [textfile collection mechanism](https://github.com/prometheus/node_exporter#textfile-collector).
In order to make use of this, one has to set up `node-exporter` according to the documentation and set the `output-file`
of this exporter to any name ending in `.prom` within said folder (and of course also `mode` to `file`).

A common use case for this is when the calculation of metrics takes particularly long and therefore can only be done
once in a while. To automate the periodic update of the output file, simply set up a cronjob.

## Example systemd unit file

```
[Unit]
Description=Prometheus disk usage exporter
Documentation=https://github.com/dundee/disk_usage_exporter

[Service]
Restart=always
User=prometheus
ExecStart=/usr/bin/disk_usage_exporter $ARGS
ExecReload=/bin/kill -HUP $MAINPID
TimeoutStopSec=20s
SendSIGKILL=no

[Install]
WantedBy=multi-user.target
```
