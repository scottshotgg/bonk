k logs ingress-nginx-controller-cf668668c-hmz4f \
| grep "remote" \
| grep -v '^{"remote":{"addr":"10.32'