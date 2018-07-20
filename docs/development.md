# Development

*All examples below assume minikube environment.*

#### Preparations

Create namespace:

    export NAMESPACE="chronologist-dev"
    kubectl create ns ${NAMESPACE}

#### Run Grafana

Deploy Grafana `v4.6.3`:

    # Note the last chart version for `4.6.3`  is `0.8.4`
    # See: https://github.com/kubernetes/charts/tree/53d1cd54f0b710c402dfd25278a66735eba969f1/stable/grafana
    
    helm install stable/grafana --version 0.8.4 \
        --wait --debug \
        --name grafana \
        --namespace ${NAMESPACE} \
        --set server.persistentVolume.enabled=false

Enable port-forwarding for Grafana pod:

     export POD_NAME=$(kubectl get pods --namespace ${NAMESPACE} -l "app=grafana-grafana,component=grafana" -o jsonpath="{.items[0].metadata.name}")
     kubectl --namespace ${NAMESPACE} port-forward $POD_NAME 3000

Export Grafana variables, getting password for user `admin`:

    export GRAFANA_ADDR="http://localhost:3000"
    export GRAFANA_PASSWORD=$(kubectl get secret --namespace ${NAMESPACE} grafana -o jsonpath="{.data.grafana-admin-password}" | base64 --decode ; echo)

Create API key:

    export GRAFANA_API_KEY=$(curl -sS -XPOST "${GRAFANA_ADDR}/api/auth/keys" \
        --user "admin:${GRAFANA_PASSWORD}" \
        -H "Content-Type: application/json" \
        -d '{"name": "chronologist", "role": "Editor"}' \
        | jq -r ".key")

Put `CHRONOLOGIST_GRAFANA_ADDR` and `CHRONOLOGIST_GRAFANA_API_KEY` in your `.env` file
to make Chronologist use that Grafana when running locally:

    cat<<EOF > .env
    CHRONOLOGIST_GRAFANA_ADDR=$GRAFANA_ADDR
    CHRONOLOGIST_GRAFANA_API_KEY=$GRAFANA_API_KEY
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

#### Run Chronologist in Minikube

Now you can shutdown you local Chronologist and deploy it to Minikube.

Build Docker image:

    docker image build -t hypnoglow/chronologist:dirty .

Push image to Minikube:

    docker save hypnoglow/chronologist:dirty | (eval $(minikube docker-env) && docker load)

Deploy Chronologist:

    helm upgrade chronologist ./deployment/chart/chronologist \
        --install --namespace ${NAMESPACE} --wait --debug \
        --set image.tag="dirty" \
        --set grafana.addr="http://grafana" \
        --set grafana.apiKey=${GRAFANA_API_KEY}

Chronologist is ready!

Refer to "Make it work!" section above to deploy some release again for
testing purposes.

#### Cleanup

    helm delete --purge chronologist grafana
