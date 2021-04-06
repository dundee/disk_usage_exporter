# Disk Usage Prometheus Exporter

Provides detailed info about disk usage of the selected filesystem path.

## Usage

```
Usage:
  disk_usage_exporter [flags]

Flags:
  -t, --analyze-interval int   How often the path should be analyzed (in seconds, detaults to 5 minutes) (default 300)
  -p, --analyzed-path string   Path where to analyze disk usage (default "/")
  -b, --bind-address string    Address to bind to (default "0.0.0.0:9108")
  -c, --config string          config file (default is $HOME/.disk_usage_exporter.yaml)
  -l, --dir-level int          Directory nesting level to show (0 = only selected dir) (default 1)
  -h, --help                   help for disk_usage_exporter
  -i, --ignore-dirs strings    Absolute paths to ignore (separated by comma) (default [/proc,/dev,/sys,/run])
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
```

## Example Prometheus queries

Disk usage of `/var` directory:

```
sum(node_disk_usage_bytes{path=~"/var.*"})
```

## Example config file

`~/.disk_usage_exporter.yaml`:
```yaml
analyze-interval: 300
analyzed-path: /
bind-address: 0.0.0.0:9995
dir-level: 2
ignore-dirs:
- /proc
- /dev
- /sys
- /run
```

## Prometheus scrape config

```yaml
scrape_configs:
  - job_name: 'disk-usage'
    static_configs:
    - targets: ['localhost:9995']
```