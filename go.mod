module github.com/open-cluster-management/clusterclaims-controller

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/openshift/hive/apis v0.0.0-20210624144808-697460baf215
	go.uber.org/zap v1.17.0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
)
