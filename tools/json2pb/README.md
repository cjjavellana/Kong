# JSON to ProtoBuf Schema

A Utility tool for generating protobuf schema given a JSON message.

## Note

This does not generate a completely accurate ProtoBuf schema yet but for the most part
generates well enough to correctly generate 80%. The remaining 20% is left up to you to 
nail down the detail.

Useful for generating large schema sets like those used by [Kong](https://docs.konghq.com/) API
Gateway. 