#!/bin/bash

export GOLANGCI_VERSION=${GOLANGCI_VERSION:-'2.1.6'}

curl -LO https://github.com/golangci/golangci-lint/raw/refs/tags/v${GOLANGCI_VERSION}/.golangci.reference.yml

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_VERSION}

# -

cat <<EOF >./header.golangci.yml
# --- 
#
# golangci-lint:       https://golangci-lint.run/
# false-positives:     https://golangci-lint.run/usage/false-positives/
# 
# -
EOF

cat ./header.golangci.yml | tee ./.golangci.yaml
cat ./.golangci.reference.yml | tee -a ./.golangci.yaml

rm ./header.golangci.yml
rm ./.golangci.reference.yml