kubectl delete -f ./config/tekton/task-source-to-image.yaml
kubectl delete -f ./config/tekton/task-deploy-on-knative.yaml
kubectl delete -f ./config/tekton/pipeline-build-and-deploy.yaml
kubectl delete -f ./config/tekton/resource-contract-git.yaml
kubectl delete -f ./config/tekton/serviceaccount.yaml