# WorkspaceOne Prometheus Exporter

WorkspaceOne UEM Prometheus exporter

## Description

This exporter is used to export WorkspaceOne UEM value to OpenMetrics format.

## Accessing the metrics

- Default port: 9740
- Endpoint: /metrics

example: `http://localhost:9740/metrics`

## Metrics

| Metric | Description | Labels | Type | Implemented |
| ------ | ----------- | ------ | ---- | ----------- |
| `device_number` | Number of devices | `none` | Gauge | yes |
| `device_os` | Number of devices by OS | `os` | Gauge | yes |
| `device_offline` | Number of offline devices | `none` | Gauge | yes |
| `device_online` | Number of online devices | `none` | Gauge | yes |
| `device_per_tag`| Number of devices by tag | `tag` | Gauge | no |
| `device_offline_per_tag`| Number of devices by tag | `tag` | Gauge | no |
| `device_online_per_tag`| Number of devices by tag | `tag` | Gauge | no |

## Environment variables

| Variable | Description |
| -------- | ----------- |
| `WS1_AUTH_KEY` | WorkspaceOne UEM user Auth Key |
| `WS1_TENANT_KEY` | WorkspaceOne UEM tenant key |
| `WS1_URL` | WorkspaceOne UEM base API URL endpoint, must finished by /API |
| `WS1_LGID` | WorkspaceOne UEM highest Group ID |
| `WS1_INTERVAL` | Interval between each WS1 check to it's enrolled devices in minutes |
| `TAG_FILTER` | String to filter Tag by it |

## Usage

## Filter by Tag

You can filter the devices by tag by using the `TAG_FILTER` environment variable.
It will enable the system to only keep the tags that contain the string you set.

### Compile from source

```bash
go build .
```

### Running it locally

```bash
go run .
```

### Run it as a Docker container

```bash
docker build -t w1-prometheus-exporter .
docker run -d -p 9740:9740 w1-prometheus-exporter
```

## Useful links

- [WorkspaceOne UEM API Reference](https://docs.vmware.com/en/VMware-Workspace-ONE-UEM/services/UEM_ConsoleBasics/GUID-BF20C949-5065-4DCF-889D-1E0151016B5A.html)
- [WorkspaceOne UEM API Explorer](https://as1506.awmdm.com/api/help/)
- [WorkspaceOne API doc pdf (UEM 9.1), is old but still has more interesting details](./doc/VMware%20AirWatch%20REST%20API%20v9_1.pdf)
