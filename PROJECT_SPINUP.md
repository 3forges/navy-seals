# Spinning the Gin project

```bash
# --
# 1./ Create the go.mod
go mod init kairos.io/seals

# --
# 2./ 
go get -u github.com/gin-gonic/gin

go mod tidy
```

```bash

go get -u github.com/hashicorp/vault/api@v1.14.8
go mod tidy

# ---
# OKAY WE REALLY HAVE A BIG ISSUE HERE: 
# the vault API does not exist in the last opensource version

# $ go get -u github.com/hashicorp/vault/api@v1.14.8
# go mod tidy
# go: downloading github.com/hashicorp/vault v1.14.8
# go: module github.com/hashicorp/vault@v1.14.8 found, but does not contain package github.com/hashicorp/vault/api

go get -u github.com/hashicorp/vault@v1.14.8
go mod tidy

# --
# I have issues installing from github the golang packages from openbao
go get -u github.com/hashicorp/openbao/api
go mod tidy

```

And I add in main.go:

```Golang
package main

import (
  "fmt"
  vapi "github.com/hashicorp/vault/api"
)

```


For the logging in CLI mode (to change for oglang gin logging appropriate logging framework):

```bash
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
```


OKAY ABOut the client I found : 

```bash
go get github.com/openbao/openbao/api/v2
```