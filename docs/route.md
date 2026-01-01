# Route

This document serves as a breakdown of Dino Route.

## Core Concepts

Dino Routes are used when routing tunnel traffic after exiting the tunnel. The routing table will match based on Route `hostname` and dial the destination provided at Route creation. 

example Route creation using Dino CLI.

```bash
dino route add \
    -a localhost:8888 \ # local address
    -p whoami.dino.local \ # hostname
    -r http \ # protocol
    -x hello \ # route name
    -t api.dino.local:50051 # api server endpoint
```

