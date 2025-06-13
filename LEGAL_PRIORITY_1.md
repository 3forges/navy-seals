We have a problem witht he vault api:

* It has dependency to `github.com/hashicorp/vault`: and there is a version from which its license is non open source.
* actually, the vault api go package source is in the same repo (monorepo) than the vault : https://github.com/hashicorp/vault/tree/v1.14.8/api


```bash
$ go get -u github.com/hashicorp/vault/api v1.14.8

go mod tidy
go: downloading github.com/hashicorp/vault/api v1.20.0
go: downloading github.com/hashicorp/vault v1.19.5

```

* Last vault opensource release: <https://github.com/hashicorp/vault/blob/v1.14.8/LICENSE>
* First non opensource licensed vault: <https://github.com/hashicorp/vault/blob/v1.14.9/LICENSE>

