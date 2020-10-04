# Compiling Kong API Gateway from Source

## Directory Description

### Smitz
An in-container agent that will be deployed alongside `Kong`. Smitz provide the authorization to
interact with `Kong`'s admin services via exposed gRPC services.

### Cyclops
Is the API Gateway Management Server responsible for monitoring API Gateway status and for
propagating config and deployment changes.

### Examples
A directory of various examples.

### Kong
Contains Kong API Gateway Configurations used for development & debugging.

### PGDATA
A placeholder directory for postgres data. This directory is mounted into the postgres container 
(see docker-compose.yml)

This directory is committed to ensure that `docker-compose` will not fail due to a missing directory.

### ProtoBuf
TODO

### Tools
TODO

## Getting Started

Building
```
$ docker build -t kong:latest
```

Running
```
$ docker-compose up kong
```
