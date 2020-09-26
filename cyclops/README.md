# Cyclops

The Kong Cluster Admin Console.

Cyclops is the one-stop shop for managing a cluster of Kong Api Gateways. At a glance, it allows portal admins to develop a sense of awareness of the health of the entire Kong Api Gateway Cluster.

## How does cyclops work

At Kong Api Gateway container startup, it registers its address to Cyclops. Cyclops adds it to the pool of managed API Gateway addresses for monitoring.

Cyclops work hand-in-hand with `smitz` to query each Kong Api Gateway instance as at pre-defined interval (default: 10s).

If Cyclops does not receive a response from an Api Gateway instance for 3 successive queries, cyclops marks that Api Gateway instance as dead and asks Openshift to terminate that instance.
