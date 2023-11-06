# Kong Ingress Controller with Go

## Steps to run the demonstration:

1. Install minikube and start it

minikube start

2. Install Kong Ingress Controller on minikube: https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/minikube/
- Test wich echo-server
- Create environment variable: export PROXY_IP=$(minikube service -n kong kong-proxy --url | head -1)

3. Build the API: 
cd code
go build -o ../api.bin cmd/main.go

4. Create a Docker image: 
If using minikube and don't want to push image to a repository, then point your local Docker client to Minikube's implementation: eval $(minikube -p minikube docker-env) --- use the same shell.
cd ..
docker build -t api:v001 .

5. Deploy the application, create a service and an ingress rule: 
```
kubectl apply -f redis.yaml
kubectl apply -f serviceDeployment.yaml
kubectl apply -f ingress-rule.yaml

kubectl delete -f redis.yaml
kubectl delete -f serviceDeployment.yaml
kubectl delete -f ingress-rule.yaml
```

6. Test access using the proxy IP: 
```
curl -i $PROXY_IP/actuator/health
curl -i --header "Content-Type: application/json" --request POST --data '{"data" : "save this", "onetime" : false}' $PROXY_IP/api/note
curl -i  $PROXY_IP/api/note/b61bc30d-8b2c-41e7-8df7-36a262826f44

```

## Setting up a rate limiter on your api

7. Create a plugin: 

kubectl apply -f rate-limiter.yaml

8. Now test resource. Try more than 5 times a minute: 
```
curl -i $PROXY_IP/actuator/health
```


