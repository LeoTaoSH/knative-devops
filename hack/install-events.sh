kubectl label namespace default knative-eventing-injection=enabled
kn service create --image docker.io/daisyycguo/service-details service-details --env EMAIL=zhanggbj@cn.ibm.com
kubectl apply -f ./config/eventing/serviceaccount-service.yaml
kubectl apply -f ./config/eventing/eventsource-service.yaml
kubectl apply -f ./config/eventing/trigger-service.yaml