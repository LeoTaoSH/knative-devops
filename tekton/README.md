# Tekton samples

## Install

1. Install Pipeline and trigger
```
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
```

2. Run samples

```
cd examples
kubectl apply -f role-resources
kubectl apply -f triggertemplates/triggertemplate.yaml
kubectl apply -f triggerbindings/triggerbinding.yaml
kubectl apply -f eventlisteners/eventlistener.yaml
```

3. Set up an ingress

```

```