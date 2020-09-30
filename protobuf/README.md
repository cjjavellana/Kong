# Inter-Process Communication Schema

This directory contains protobuf message definitions used between Cyclops & Smitz. 

## Quick Start

1. Install Protobuf Compiler. [See Compiling Protocol Buffers](https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers)

2. Generating Protobuf Stubs

    ```bash
    $ protoc -I=./ -I=./include --go_out=../smitz --go-grpc_out=../smitz ./smitz.proto
    ```