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

#### Option A: Kubernetes manifests

1. Clone this repo
2. Navigate to [deployment/manifests](deployment/manifests)
3. Review and change [configmap.yaml](deployment/manifests/configmap.yaml) to set desired configuration.
    1. Review and change [rbac.yaml](deployment/manifests/rbac.yaml), at least to set up desired namespace.
    1. Optionally review and change [deployment.yaml](deployment/manifests/deployment.yaml)
4. Deploy to Kubernetes cluster:

        kubectl apply -f ./

#### Option B: Helm Chart

1. Clone this repo
2. Navigate to [deployment/chart](deployment/chart)
3. Review [values.yaml](deployment/chart/chronologist/values.yaml)
    1. Create a new values-file that reflects your desired configuration and pass
    it to chart installation command below,
    2. OR you can pass values using `--set` flag when installing chart.
4. Deploy chart to Kubernetes cluster (passing values from step 3):

        helm install ./chronologist --name chronologist \
            --values ./my-values.yaml

## Contributing

Contributions are welcome!

See [docs/development.md](docs/development.md) for detailed instructions on 
how to run development environment for Chronologist.
