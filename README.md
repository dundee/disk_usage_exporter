# Disk Usage Prometheus Exporter

Provides detailed info about disk usage of the selected filesystem path.

## Usage

```
Usage:
  gdu-prometheus-exporter [flags]

Flags:
  -p, --analyzed-path string   Path where to analyze disk usage (default "/")
  -b, --bind-address string    Address to bind to (default "0.0.0.0:9108")
  -c, --config string          config file (default is $HOME/.gdu-prometheus-exporter.yaml)
  -l, --dir-level int          Directory nesting level to show (0 = only selected dir) (default 1)
  -h, --help                   help for gdu-prometheus-exporter
  -i, --ignore-dirs strings    Absolute paths to ignore (separated by comma) (default [/proc,/dev,/sys,/run])
```

## Example output

```
# HELP gdu_disk_usage Disk usage of the directory
# TYPE gdu_disk_usage gauge
gdu_disk_usage{directory="/",level="0"} 1.3118748672e+11
gdu_disk_usage{directory="/bin",level="1"} 0
gdu_disk_usage{directory="/boot",level="1"} 1.29736704e+08
gdu_disk_usage{directory="/etc",level="1"} 1.3078528e+07
gdu_disk_usage{directory="/home",level="1"} 8.5037293568e+10
gdu_disk_usage{directory="/lib",level="1"} 0
gdu_disk_usage{directory="/lib64",level="1"} 0
gdu_disk_usage{directory="/lost+found",level="1"} 4096
gdu_disk_usage{directory="/mnt",level="1"} 4096
gdu_disk_usage{directory="/opt",level="1"} 2.979229696e+09
gdu_disk_usage{directory="/root",level="1"} 4096
gdu_disk_usage{directory="/sbin",level="1"} 0
gdu_disk_usage{directory="/snap",level="1"} 0
gdu_disk_usage{directory="/srv",level="1"} 4.988928e+06
gdu_disk_usage{directory="/tmp",level="1"} 1.3561856e+07
gdu_disk_usage{directory="/usr",level="1"} 1.7940340736e+10
gdu_disk_usage{directory="/var",level="1"} 2.506924032e+10
```

## Example config file

`~/.gdu-prometheus-exporter.yaml`:
```yaml
analyzed-path: /
bind-address: 0.0.0.0:9108
dir-level: 1
ignore-dirs:
- /proc
- /dev
- /sys
- /run
```