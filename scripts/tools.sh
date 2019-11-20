#!/bin/sh
#
# scripts/tools.sh
#
# Installs tools dependencies
#

go install \
  github.com/golangci/golangci-lint/cmd/golangci-lint \
  github.com/maxbrunsfeld/counterfeiter/v6
