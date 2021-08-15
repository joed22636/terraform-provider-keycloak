# terraform-provider-keycloak
Terraform provider for [Keycloak](https://www.keycloak.org/).

[![CircleCI](https://circleci.com/gh/joed22636/terraform-provider-keycloak.svg?style=shield)](https://circleci.com/gh/joed22636/terraform-provider-keycloak)

# Installation

v1.0.3 and above can be installed automatically using Terraform 0.13 by using the `terraform` configuration block:

```hcl
terraform {
  required_providers {
    keycloak = {
      source = "joed22636/keycloak"
      version = ">= 1.0.3"
    }
  }
}
```
# Supported Versions

The following versions are used when running acceptance tests in CI:

- 11.0.3

# Releases

This provider uses [GoReleaser](https://goreleaser.com/) to build and publish releases. Each release published to GitHub
contains binary files for Linux, macOS (darwin), and Windows, as configured within the [`.goreleaser.yml`](https://github.com/joed22636/terraform-provider-keycloak/blob/master/.goreleaser.yml)
file.

Each release also contains a `terraform-provider-keycloak_${RELEASE_VERSION}_SHA256SUMS` file, accompanied by a signature
created by a PGP key with the fingerprint `C508 6791 5E11 6CD2`. This key can be found on my Keybase account at https://keybase.io/joed22636.

You can find the list of releases [here](https://github.com/joed22636/terraform-provider-keycloak/releases).
<!-- You can find the changelog for each version [here](https://github.com/joed22636/terraform-provider-keycloak/blob/master/CHANGELOG.md). -->

# Development

This project requires Go 1.15 and Terraform 0.13.
This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management, which allows this project to exist outside of an existing GOPATH.

After cloning the repository, you can build the project by running `make build`.

## Local Environment

You can spin up a local developer environment via [Docker Compose](https://docs.docker.com/compose/) by running `make local`.
This will spin up a few containers for Keycloak, PostgreSQL, and OpenLDAP, which can be used for testing the provider.
This environment and its setup via `make local` is not intended for production use.

Note: The setup scripts require the [jq](https://stedolan.github.io/jq/) command line utility.

### Macos: 
```
brew install jq
brew install terraform@0.13
```

## Tests

use keycloak/keycloak_client_test.go to debug http rest apis

Every resource supported by this provider will have a reasonable amount of acceptance test coverage.

```
make local
```

You can run acceptance tests against a Keycloak instance by running `make testacc`. You will need to supply some environment
variables in order to set up the provider during tests. Here is an example for running tests against a local environment
that was created via `make local`:

```
KEYCLOAK_CLIENT_ID=terraform \
KEYCLOAK_CLIENT_SECRET=884e0f95-0f42-4a63-9b1f-94274655669e \
KEYCLOAK_CLIENT_TIMEOUT=5 \
KEYCLOAK_REALM=master \
KEYCLOAK_URL="http://localhost:8080" \
make test > test.log
```


```
make testacc > testacc.log
```

Result strings:
```
--- PASS
--- FAIL
```

Run specific test from CLI, e.g.:
```
export TF_LOG=DEBUG                                                                    
go test -v -run <name> <package>
go test -v -run TestAccKeycloakCustomUserFederation_createAfterManualDestroy ./provider
go test -v -run TestAccKeycloakUsersPermission_basic ./provider
```
## TODO

* credential handling 
  * smtp (not yet needed)
  * ldap
  * idp 
* import of master realm security defenses does not work
* import of enabled event types does not work
* import missing
    client scope - include in token flag
    bearer only - standard flow enabled is set to true at import 
* [not yet needed] - webUI > client scope > [select one] > scope > realm/client role mappings are not managable afais
* [not yet needed] - idp > some items (allowedClockSkew,forwardParameters,prompt selection)
* [not yet needed] - webUI > realm > token > some settings (ssoSessionIdleTimeoutRememberMe, ssoSessionMaxLifespanRememberMe)
* [not yet needed] - webUI > cilents > select one > Permissions
* [not yet needed] - webUI > cilents > Fine Grain OpenID Connect Configuration 
* [not yet needed] - webUI > cilents > OAuth 2.0 Mutual TLS Certificate Bound Access Tokens Enabled 
* [not yet needed] - webUI > cilents > consent related settings
* [not yet needed] - webUI > cilents > authorization enabled and related settings
* [not yet needed] - webUI > user federation > some settings (enabled tls, debug, Enable the LDAPv3 Password Modify Extended Operation)
* cicleci reference in readme - badge - CI build
* built in flow property is not kept

## Extend - classes - architecture

see bindings resource, connection pooling

### HTTP API - json model
 JSON model and API CRUD methods, package: keycloak
  
 HTTP response <--> JSON model
 file: e.g. authentication_bindings.go

### Terraform model 
Terraform schema model and CRUD methods, starts with , package: provider
TF scheme  <--> JSON model
file: resource_*, e.g. resource_keycloak_authentication_bindings.go

# License

[MIT](https://github.com/joed22636/terraform-provider-keycloak/blob/master/LICENSE)
