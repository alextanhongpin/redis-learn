## Kubectl 

Install latest kubectl (>1.14) that supports `kustomization`:

```bash
$ brew install kubernetes-cli

# Create the customization file.
$ k apply -k .

# Deploy.
$ k apply -f redis.yaml

# Find the port 
$ k get svc
NAME            TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
kubernetes      ClusterIP   10.96.0.1      <none>        443/TCP          30m
redis-service   NodePort    10.107.237.0   <none>        6379:31344/TCP   2m

# Replace the port in main.go.
```

 
