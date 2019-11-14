kubectl apply -f ./config/tekton/task-source-to-image.yaml
kubectl apply -f ./config/tekton/task-deploy-on-knative.yaml
kubectl apply -f ./config/tekton/pipeline-build-and-deploy.yaml
kubectl apply -f ./config/tekton/resource-contract-git.yaml
kubectl apply -f ./config/tekton/serviceaccount.yaml