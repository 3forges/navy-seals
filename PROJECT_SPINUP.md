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

For the logging in CLI mode (to change for Golang gin logging appropriate logging framework):

```bash
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
```

OKAY About the client I found :

```bash
go get github.com/openbao/openbao/api/v2
```