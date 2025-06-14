# Seals

I have now a very simple executable that can seal/unseal a vault:
* deployed in dev mode
* with only one single unseal key

So now to go further in tests, and that before coding anything like the REST API using gin, I need:
* A new helm deployment of the vault
* with a csi driver and a dynamic volume provisioner: 
  * that is required for the vault to be able to persist its data
  * I will also have to configure values.yml for the vault ot use the storage class.
* I will have to configure the helm chart with a custom init process to be able to set the threshold / number of keys:
  * like with a `` command
  * the created unseal keys must then be stored somewhere my rest api can access them once
  * note that I could check the unseal keys can be encrypted with GPG Keys for example GPG KEys coming from keybase.io as far as I could understand (maybe we will need our own private GPG Keys service)
  * And here is below the part of the [`values.yaml`](https://openbao.github.io/openbao-helm/charts/openbao/values.yaml) that I have to work with to customize the vault init process I think:

```Yaml
# This is from https://openbao.github.io/openbao-helm/charts/openbao/values.yaml
# -- 

  # Used to define commands to run after the pod is ready.
  # This can be used to automate processes such as initialization
  # or boostrapping auth methods.
  postStart: []
  # - /bin/sh
  # - -c
  # - /vault/userconfig/myscript/run.sh
```


Oh to initializes and all here are examples I found too on github:

* https://github.com/linode/docs/blob/e6fb945938faa63b89625f832f2162732a053935/docs/guides/security/secrets-management/deploy-openbao-on-linode-kubernetes-engine/index.md#initialize-and-unseal-the-openbao-development-server

* also : 
  * https://github.com/Alfred-Sabitzer/microk8s-ubuntu/tree/ea12edf7956130b066263024d35605eee49c7ebd/setup/openBao
  * https://github.com/Alfred-Sabitzer/microk8s-ubuntu/blob/ea12edf7956130b066263024d35605eee49c7ebd/setup/openBao/openBao.sh#L37

## References

* https://github.com/lrstanley/vault-unseal
* https://github.com/hashicorp/hello-vault-go/tree/main/sample-app : here very good examples.


## ANNEX: Testing

So qickly provision an openBAO in a Kind Cluster like this:

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


# helm delete ${HELM_RELEASE_NAME} -n ${K8S_NS}
# kubectl -n ${K8S_NS} port-forward service/pesto-openbao --address=192.168.1.16 8200:8200

helm get manifest -n ${K8S_NS} ${HELM_RELEASE_NAME}
helm status -n ${K8S_NS} ${HELM_RELEASE_NAME}

kubectl --context kind-openbao-cluster get all -n pesto-openbao


```

Now the Vault is not initialiazed, and you need to initialize it like this:

```bash
export BAO_ADDR='http://192.168.1.16:8200'
export BAO_TOKEN='root'

bao operator init -key-shares=17 -key-threshold=7



# $ bao operator init -key-shares=17 -key-threshold=7
# Unseal Key 1: CcsoQoRcdixD/CbOiGX96EI1pc43iUe7eb6EVd1E8NYv
# Unseal Key 2: IygQ4+7Jpig3Pf+S15ZuGtbje1wri6MoYftrpfE6tu9h
# Unseal Key 3: VLxqjpE/tCKjlGFIeqLJ8c/+oWoWK/a2BTbvDZE+xmb6
# Unseal Key 4: ZXs016hHR8stHVhZ+c5+eTZdYYv6bVPEsvlEl2kpAjCH
# Unseal Key 5: ZpuWEez/hDSbsTkkjRivqhPjtUA0eu9p3i9dWkHeOcpK
# Unseal Key 6: vYEBg2CH23dW0InRtThATO8T9lLPKWUBuYk/RudycLiw
# Unseal Key 7: qq3clb3if/DgFCzZnG52bxbm5fw/a6bOajDjDMVMLsL5
# Unseal Key 8: RfxYkedQQJTaK58LpbiaOgh7Kw1vOcMlyd3l3WPU//5m
# Unseal Key 9: YR+erJ+pksWZeTvaRLZeRS4EQN3C8oSEWOSMfq+5KKtz
# Unseal Key 10: ucDQj6u2I5Hs0Kn+lrqIRziGvwzsbF+PXj0oGwzf9PYX
# Unseal Key 11: fVujsysYPEb74ihH4kE702/LOhwqYjpjawLnvnDX+aAf
# Unseal Key 12: SFH5CJINrsf9TbovF17c+oL3OMR8OkkN0DN8F09c48xX
# Unseal Key 13: lhrIKkx8PS0GYA7sykhoyPwtKOURBR2OO3hURN0PuoZa
# Unseal Key 14: IoXSLjPb8uMh0fYPCQjMpw8wQmOMRlh9jhCxCjiHPAcl
# Unseal Key 15: UqYFPMIwxs0lRSopsC2Sg8skPw0MVNBI/iMskD2phw1t
# Unseal Key 16: mndJ0iCKBQ4rvnhPZlZga/pnRua5vF68SddSOugxpwtM
# Unseal Key 17: WBmw2edH5mRU9B4nzJdjLATsY4OlTSFsT7Wf3SlbDxMT

# Initial Root Token: s.hWyIRGkJ59uggT6lSyHBeY3b

# Vault initialized with 17 key shares and a key threshold of 7. Please securely distribute the key shares printed above. When the Vault is re-sealed, restarted, or stopped, you must supply at least 7 of these keys to unseal it before it can start servicing requests.

# Vault does not store the generated root key. Without at least 7 keys to reconstruct the root key, Vault will remain permanently sealed!
# 
# It is possible to generate new unseal keys, provided you have a quorum of
# existing unseal keys shares. See "bao operator rekey" for more information.


```

When you need to delete the vault and its persistent volumes so restart from a fresh non initialized vault:

```bash
export HELM_RELEASE_NAME=${HELM_RELEASE_NAME:-'pesto-openbao'}
export K8S_NS=${K8S_NS:-'pesto-openbao'}

helm delete ${HELM_RELEASE_NAME} -n ${K8S_NS}

kubectl -n ${K8S_NS} delete persistentvolumeclaim/data-${HELM_RELEASE_NAME}-0


```

* Now you can run navy seals as a CLI like that:

```bash

export UNSEAL_TOKENS='CcsoQoRcdixD/CbOiGX96EI1pc43iUe7eb6EVd1E8NYv,IygQ4+7Jpig3Pf+S15ZuGtbje1wri6MoYftrpfE6tu9h,VLxqjpE/tCKjlGFIeqLJ8c/+oWoWK/a2BTbvDZE+xmb6,ZXs016hHR8stHVhZ+c5+eTZdYYv6bVPEsvlEl2kpAjCH,ZpuWEez/hDSbsTkkjRivqhPjtUA0eu9p3i9dWkHeOcpK,vYEBg2CH23dW0InRtThATO8T9lLPKWUBuYk/RudycLiw,qq3clb3if/DgFCzZnG52bxbm5fw/a6bOajDjDMVMLsL5,RfxYkedQQJTaK58LpbiaOgh7Kw1vOcMlyd3l3WPU//5m,YR+erJ+pksWZeTvaRLZeRS4EQN3C8oSEWOSMfq+5KKtz,ucDQj6u2I5Hs0Kn+lrqIRziGvwzsbF+PXj0oGwzf9PYX,fVujsysYPEb74ihH4kE702/LOhwqYjpjawLnvnDX+aAf,SFH5CJINrsf9TbovF17c+oL3OMR8OkkN0DN8F09c48xX,lhrIKkx8PS0GYA7sykhoyPwtKOURBR2OO3hURN0PuoZa,IoXSLjPb8uMh0fYPCQjMpw8wQmOMRlh9jhCxCjiHPAcl,UqYFPMIwxs0lRSopsC2Sg8skPw0MVNBI/iMskD2phw1t,mndJ0iCKBQ4rvnhPZlZga/pnRua5vF68SddSOugxpwtM,WBmw2edH5mRU9B4nzJdjLATsY4OlTSFsT7Wf3SlbDxMT'


echo "UNSEAL_TOKENS='${UNSEAL_TOKENS}'" > ./.env

unset UNSEAL_TOKENS

# Always make sure that UNSEAL_TOKENS is unset and only the vale set i n [./.env] file is ruling

# ./dist/bin/navy-seals -unseal-keys-nb 23 -unseal-keys-treshold 11
make
./dist/bin/navy-seals

```