# run-server

Small HTTPS server meant for testing in local environment.

## Dependencies

- [mkcert](https://github.com/FiloSottile/mkcert/?tab=readme-ov-file#installation)
- Golang

## Setup

The registry must run on HTTPS for OpenTofu to be able to communicate with it.
We use `mkcert` to create a locally-trusted SSL certificate.

```sh
# Create and install CA in system trust store.
mkcert -install
# Create a new certificate valid for localhost
mkcert -cert-file cmd/run-server/localhost.pem -key-file cmd/run-server/localhost-key.pem localhost
```

## Run

First, populate the `generated` folder with the registry content

```sh
go run ./cmd/generate-v1 --destination ../generated
```

Then you can start the server

```sh
go run ./cmd/run-server/ -certificate cmd/run-server/localhost.pem -key cmd/run-server/localhost-key.pem
```

## Test

Check if the registry is reachable with `curl`

```sh
$ curl -D - https://localhost:8443/.well-known/terraform.json
HTTP/2 200
accept-ranges: bytes
content-type: application/json
last-modified: Mon, 11 Mar 2024 14:31:37 GMT
content-length: 72
date: Mon, 11 Mar 2024 14:34:19 GMT

{
          "modules.v1": "/v1/modules/",
          "providers.v1": "/v1/providers/"
}
```

Fetch providers

```sh
terraform {
  required_providers {
    random = {
      source  = "localhost:8443/hashicorp/random"
      version = "~> 3"
    }
  }
}
```

Fetch modules

```sh
module "consul" {
  source = "localhost.:8443/hashicorp/consul/aws" # Keep the dot after localhost, it is not a typo
  version = "0.11.0"
}
```
