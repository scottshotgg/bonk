# Bonk

`bonk` is currently used in my homelab to propagate iptables bans based on L7 information. To explain further, the router accepts a request, shoots it over to the k8s cluster, where it is then disected and routed to the appropriate service for processing. In that routing, there could be certain tell-tale sign that this is fraudulent request for a number of different reason which are detailed below and codified in `pkg/engine`:
- spoofed IP address
- invalid path
- invalid URL
- using my external IP directly instead of my URL
- No HTTP method, host, or URI
- Specific paths and user-agents that are known to not exist (i.e, `cgi-bin`, wordpress stuff, etc) 