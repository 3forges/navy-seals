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
