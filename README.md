# Chronologist 🎞

[![CircleCI](https://circleci.com/gh/hypnoglow/chronologist.svg?style=shield)](https://circleci.com/gh/hypnoglow/chronologist)

Chronologist is a Kubernetes controller that syncs your Helm chart deployments 
with Grafana annotations.

![screenshot](/screenshot.png)

Key features:

- For each Helm release you install/upgrade creates related Grafana annotation
- Determines if it was a rollout or rollback and tags annotation appropriately
- When you purge delete a release, deletes corresponding annotation

## Deployment

The easiest way to install Chronologist into your Kubernetes cluster is to use
Helm chart.

    helm repo add chronologist https://hypnoglow.github.io/chronologist
    helm install chronologist/chronologist

## Contributing

Contributions are welcome!

See [docs/development.md](docs/development.md) for detailed instructions on 
how to run development environment for Chronologist.
