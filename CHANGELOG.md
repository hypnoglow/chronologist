# Change log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.2.0]

### Added

- Add abbility to watch secrets.

    Helm uses Kubernetes ConfigMaps as a backend for storing releases.
    However, [since Helm 2.7](https://github.com/helm/helm/pull/2721),
    it is possible to enable Kubernetes Secrets as a backend.

    Chronologist now supports secret backend; thus it can be configured
    either to watch ConfigMaps or Secrets.

### Changed

- **BREAKING** Revamped chart configuration, making it simpler to configure via `values.yaml`.

### Security

- Chronologist container now runs as non-root.

## [0.1.1] - 2018-06-13

### Added

- Introduce kubernetes secret in the chart to store sensitive data
like Grafana API key.

## [0.1.0] - 2018-06-09

Initial release.
