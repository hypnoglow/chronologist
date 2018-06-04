# Development

*All examples below assume minikube environment.*

#### Preparations

Create namespace:

    kubectl create ns chronologist-dev

#### Run Grafana

Add this line to `/etc/hosts` if you want to access an ingress in minikube:

    192.168.99.100 grafana.chronologist.minikube.local

Deploy Grafana `v4.6.3`:

    # Note the last chart version for `4.6.3`  is `0.8.4`
    # See: https://github.com/kubernetes/charts/tree/53d1cd54f0b710c402dfd25278a66735eba969f1/stable/grafana
    
    helm install stable/grafana --version 0.8.4 \
        --wait --debug \
        --name grafana-chronologist \
        --namespace chronologist-dev \
        --set server.persistentVolume.enabled=false \
        --set server.ingress.enabled=true \
        --set server.ingress.hosts.0=grafana.chronologist.minikube.local
    
    export GRAFANA_ADDR="http://grafana.chronologist.minikube.local"

Get password for user `admin`:

    export GRAFANA_PASSWORD=$(kubectl get secret --namespace chronologist-dev grafana-chronologist -o jsonpath="{.data.grafana-admin-password}" | base64 --decode ; echo)

Create API key:

    export GRAFANA_API_KEY=$(curl -sS -XPOST "${GRAFANA_ADDR}/api/auth/keys" \
        --user "admin:${GRAFANA_PASSWORD}" \
        -H "Content-Type: application/json" \
        -d '{"name": "chronologist", "role": "Editor"}' \
        | jq -r ".key")

Put `GRAFANA_ADDR` and `GRAFANA_API_KEY` in your `.env` file
to make Chronologist use that grafana when running locally:

    cat<<EOF > .env
    GRAFANA_ADDR=$GRAFANA_ADDR
    GRAFANA_API_KEY=$GRAFANA_API_KEY
    EOF

#### Run Chronologist locally

Build Chronologist:

    make build

Run:

    ./bin/chronologist

#### Make it work!

Deploy some helm chart:

    helm install stable/kube-ops-view --name foo

Watch Chronologist output as it creates an annotation in Grafana.

Check that annotation:
    
    curl -sS -XGET "${GRAFANA_ADDR}/api/annotations" \
        -H "Authorization: Bearer ${GRAFANA_API_KEY}"

Try delete previously deployed release and watch Chronologist deletes the annotation:

    helm delete --purge foo

Check that annotation does not exist anymore:
    
    curl -sS -XGET "${GRAFANA_ADDR}/api/annotations" \
        -H "Authorization: Bearer ${GRAFANA_API_KEY}"
