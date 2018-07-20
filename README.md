# Chronologist 🎞

[![CircleCI](https://circleci.com/gh/hypnoglow/chronologist.svg?style=shield)](https://circleci.com/gh/hypnoglow/chronologist)

Chronologist is a Kubernetes controller that syncs your Helm chart deployments 
with Grafana annotations.

![screenshot](/screenshot.png)

Key features:

- For each Helm release you install/upgrade creates related Grafana annotation
- Annotations are tagged with related info such as release name, release namespace, etc
- When you purge delete a release, deletes corresponding annotation

## Deployment

The easiest way to install Chronologist into your Kubernetes cluster is to use
Helm chart.

*Replace values below with your actual Grafana address and API key*

    helm repo add chronologist https://hypnoglow.github.io/chronologist
    helm install chronologist/chronologist \
        --set grafana.addr=http://grafana.example.com \
        --set grafana.apiKey=ABCDEF1234567890

See [values.yaml](../deployment/chart/chronologist/values.yaml) for the full list
of possible options.

## Contributing

Contributions are welcome!

See [docs/development.md](docs/development.md) for detailed instructions on 
how to run development environment for Chronologist.
