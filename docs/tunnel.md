# 🛤️ Tunnel

This document explains how Dino maps incoming tunnel traffic to internal destinations. 

## Core Concept

Dino operates a control plane for network tunnels. Functionality is different from wireguard as it does not interact with physical devices. Instead it manages the routing table that dictates where traffic goes once it exits the tunnel.

example tunnel creation using Dino CLi.
```bash
    dino tunnel add hello -t api.dino.local:50051
```

## Authentication



