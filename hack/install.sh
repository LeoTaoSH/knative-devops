kubectl label namespace default knative-eventing-injection=enabled
kn service create --image docker.io/daisyycguo/service-details service-details --env EMAIL=zhanggbj@cn.ibm.com
kubectl apply -f ./config/serviceaccount-service.yaml
kubectl apply -f ./config/eventsource-service.yaml
kubectl apply -f ./config/trigger-service.yaml