kn service delete service-details
kubectl delete -f ./config/trigger-service.yaml
kubectl delete -f ./config/eventsource-service.yaml
kubectl delete -f ./config/serviceaccount-service.yaml
kubectl label namespace default knative-eventing-injection-
kubectl delete broker default
