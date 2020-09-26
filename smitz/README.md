# SMITZ

The agent that lives within the kong container. For security reasons, kong admin port is only restricted to 127.0.0.1. Smitz apply PKI authentication to incoming requests to the admin port to authenticate the request before forwarding to the loop back address.


