# Terraform Provider for BPK.io

This Terraform Provider allows you to manage BPK.io resources through Terraform. It provides a convenient way to automate the creation, management, and deletion of your BPK.io infrastructure as part of your Infrastructure as Code workflow.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/bashou/terraform-provider-bpkio` to your Terraform provider:

```shell
go get github.com/bashou/terraform-provider-bpkio
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

### Provider Configuration

Configure the provider with your BPK.io credentials:

```hcl
provider "bpkio" {
  api_key = var.bpkio_api_key  # Can also be set via  BPKIO_API_KEY environment variable
}
```

### Resources

### Data Sources

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
