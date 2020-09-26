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

## APIs

All Kong 2.1.x admin APIs are supported. See [Admin API](https://docs.konghq.com/2.1.x/admin-api/)


