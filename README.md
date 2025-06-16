# Seals

In this first release, I iwll deliver a very simple CLI which is able:
* to init an OpenBAO vault
* to display the status
* to seal the OpenBAO vault
* to unseal the OpenBAO vault

Note that the generated unseal keys are turned into QR codes inside the `.tofu_secrets` folder.

ALso note that navy seal unseals the vault by ready the unseal keys from the same QR codes.

## Test it

You will:

* provision an openBOA vault inside a simple Kubernetes Cluster (use `kind` for example) 
* then you will run navy-seals commands
* then you will tear down the OpenBAO vault, including its data (persistent volumes), such that the vault comes back to its initial state, not initialized.

* Provision the OpenBAO Vault:

```bash
helm repo add openbao https://openbao.github.io/openbao-helm

helm search repo openbao/openbao -l

export DESIRED_CHART_VERSION=${DESIRED_CHART_VERSION:-'0.12.0'}

export HELM_KUBECONTEXT=${HELM_KUBECONTEXT:-'kind-openbao-cluster'}


export HELM_RELEASE_NAME=${HELM_RELEASE_NAME:-'pesto-openbao'}
export K8S_NS=${K8S_NS:-'pesto-openbao'}


helm install ${HELM_RELEASE_NAME} openbao/openbao --version ${DESIRED_CHART_VERSION} \
     --namespace ${K8S_NS} \
     --create-namespace \
     --set server.dev.enabled=false

```

Now the Vault is not initialiazed, and ready for the `navy-seals` tests. You can run navy seals as a CLI like that:

```bash

# build id from source

make

#./dist/bin/navy-seal -b 0.0.0.0 -p 8751
./dist/bin/navy-seal -b localhost -p 8751

# Test Using the list albums endpoint:

curl -X GET \
    --header "Content-Type: application/json" \
    http://localhost:8751/albums

curl -X GET \
    --header "Content-Type: application/json" \
    http://localhost:8751/albums | jq .

# Test adding an Album:

curl http://localhost:8751/albums \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'

curl -X GET \
    --header "Content-Type: application/json" \
    http://localhost:8751/albums | jq .

curl -X GET \
    --header "Content-Type: application/json" \
    http://localhost:8751/albums/3 | jq .

curl -X GET \
    --header "Content-Type: application/json" \
    http://localhost:8751/albums/4 | jq .

```

And when you need to delete the vault and its persistent volumes so restart from a fresh non initialized vault:

```bash
export HELM_RELEASE_NAME=${HELM_RELEASE_NAME:-'pesto-openbao'}
export K8S_NS=${K8S_NS:-'pesto-openbao'}

helm delete ${HELM_RELEASE_NAME} -n ${K8S_NS}

kubectl -n ${K8S_NS} delete persistentvolumeclaim/data-${HELM_RELEASE_NAME}-0

# kubectl -n pesto-openbao delete persistentvolumeclaim/data-pesto-openbao-0

```


## ANNEX: Gin References

* Interesting: https://mcorbin.fr/posts/2022-06-13-gin-golang/