module github.com/daisy-ycguo/knative-devops/service

go 1.12


require (
	github.com/cloudevents/sdk-go v0.10.0
	k8s.io/api kubernetes-1.15.3
	k8s.io/apimachinery kubernetes-1.15.3
	k8s.io/client-go kubernetes-1.15.3
	knative.dev/serving v0.9.0
)