# WorkspaceOne Prometheus Exporter

WorkspaceOne UEM Prometheus exporter

## Description

This exporter is used to export WorkspaceOne UEM value to OpenMetrics format.

## Accessing the metrics

- Default port: 9740
- Endpoint: /metrics

example: `http://localhost:9740/metrics`

## Metrics

| Metric | Description | Labels | Type |
| ------ | ----------- | ------ | ---- |
| `device_number` | Number of devices | `none` | Gauge |
| `device_os` | Number of devices by OS | `os` | Gauge |
| `device_os_version` | Number of devices by OS version | `os_version` | Gauge |
| `device_model` | Number of devices by model | `model` | Gauge |

## Environment variables

| Variable | Description |
| -------- | ----------- |
| `WS1_AUTH_KEY` | WorkspaceOne UEM user Auth Key |
| `WS1_TENANT_KEY` | WorkspaceOne UEM tenant key |
| `WS1_URL` | WorkspaceOne UEM base API URL endpoint, must finished by /API |
| `WS1_LGID` | WorkspaceOne UEM highest Group ID |

## Usage

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
