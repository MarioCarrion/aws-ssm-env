# AWS Systems Manager (SSM) Environment Variables

## Motivation

Storing configuration in the environment, as in a [Twelve factor App](https://12factor.net/config), is a common use case. However there are cases where there's the need to properly secure the value of those variables, for example when exposing secrets (such as usernames and passwords) or defining API Tokens.

## What is this?

`aws-ssm-env` is a Go package that allows you load configuration values from environment variables and in cases where explicitly stated it could load the values from AWS SSM.

## How does it work?

1. Define a new struct that will hold all your configuration values, decorate it using the `ssm` tag.
1. Load your configuration file, if needed.
1. Initialize [AWS SSM](https://github.com/aws/aws-sdk-go/tree/master/service/ssm)
1. Load the configuration values.

Please look at [the example](examples/main.go) for a concrete full example.

## Development

* Install [`direnv`](https://github.com/direnv/direnv)
* Install required tools using `./script/tools.sh`
* `go generate ./...` and `go test ./...` as usual
