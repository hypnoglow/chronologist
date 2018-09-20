# Sinks

Currently, only Grafana is supported as a sink for releases.

## Grafana

Release Events are recorded as Grafana Annotations.

## Other options

Why not Prometheus? Probably because it is not the right solution
for event store.

> The Pushgateway is not an event store. While you can use Prometheus as a data source for Grafana annotations, tracking something like release events has to happen with some event-logging framework. [\[1\]](1)

[1]: https://github.com/prometheus/pushgateway/tree/5d69bdfacfac6393e2b6b8d667874b603b7b04fa#non-goals


