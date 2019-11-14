kn service delete service-details
kubectl delete -f ./config/eventing/trigger-service.yaml
kubectl delete -f ./config/eventing/eventsource-service.yaml
kubectl delete -f ./config/eventing/serviceaccount-service.yaml
kubectl label namespace default knative-eventing-injection-
kubectl delete broker default
