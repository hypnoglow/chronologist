# Deployment options

The easiest way to install Chronologist into your Kubernetes cluster is to use
Helm chart. The Helm chart repo is served from this same GitHub repo using [GitHub Pages](https://github.com/hypnoglow/chronologist/tree/gh-pages).

*Replace values below with your actual Grafana address and API key*

    helm repo add chronologist https://hypnoglow.github.io/chronologist
    helm install chronologist/chronologist \
        --set config.GRAFANA_ADDR=http://grafana.example.com \
        --set secrets.GRAFANA_API_KEY=ABCDEF1234567890

See [values.yaml](../deployment/chart/chronologist/values.yaml) for the full list
of possible options.

## Alternatives

Alternatives involve cloning this repo and manipulating source files.

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
