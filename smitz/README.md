# SMITZ

An in-container proxy that proxies request to Kong's admin port.

For security reasons, kong admin port is only restricted to 127.0.0.1. Smitz apply PKI authentication to incoming requests to the admin port to authenticate the request before forwarding to the loop back address.

## How to run Smitz

### Environment Variables

Smitz does not require too many of a parameters to start. It only needs the address of its management console aka Cyclops in order to register itself.

### `CYCLOPS_URL`

Example `docker-compose.yml`:

```yaml
version: '3.1'

services:

  kong:
    image: kong
    restart: always
    environment:
      CYCLOPS_URL: http://cyclops:8080
			... other kong env variables here
```

Example `config-map` for `k8` or `Openshift`:

```yaml
TODO (cjavellana): Add Config Map Configuration 
```

### Starting via the command line

```bash
$ smitz --cyclops-url ${CYCLOPS_URL}
```

## APIs

All Kong 2.1.x admin APIs are supported. See [Admin API](https://docs.konghq.com/2.1.x/admin-api/)

## Testing gRPC Endpoints

If you are on macOS, you can install [grpcurl](https://github.com/fullstorydev/grpcurl) via Homebrew.

```bash
$ brew install grpcurl
```

Listing out Smitz gRPC endpoints.
```bash
$ grpcurl -plaintext localhost:8686 list
grpc.reflection.v1alpha.ServerReflection
kong.proxy.Admin
```

Inspecting the detail of a gRPC method
```bash
$ grpcurl -plaintext localhost:8686 list kong.proxy.Admin
kong.proxy.Admin.GetStatus
kong.proxy.Admin.NodeInfo
```

Calling a gRPC method
```bash
$ grpcurl -plaintext localhost:8686 kong.proxy.Admin/NodeInfo
{
  "plugins": {
    "availableOnServer": {
      "pluginName": [
        "One",
        "Two"
      ]
    },
    "enabledInCluster": {
      "pluginName": [
        "Three",
        "Four"
      ]
    }
  },
  "configuration": {
    "name": [
      "One",
      "Two"
    ]
  }
}
```