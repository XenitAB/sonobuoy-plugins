module github.com/xenitab/sonobuoy-plugins

go 1.15

require (
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/emicklei/go-restful v2.9.6+incompatible // indirect
	github.com/fluxcd/kustomize-controller/api v0.5.3
	github.com/fluxcd/source-controller/api v0.5.6
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/onsi/ginkgo v1.13.0
	github.com/onsi/gomega v1.10.1
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.11.0 // indirect
	google.golang.org/genproto v0.0.0-20200619004808-3e7fca5c55db // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c // indirect
	k8s.io/api v0.19.4
	k8s.io/apiextensions-apiserver v0.19.0 // indirect
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v1.19.0
	k8s.io/component-base v0.19.0
	k8s.io/klog/v2 v2.3.0
	k8s.io/kubectl v0.19.0 // indirect
	k8s.io/kubernetes v1.19.0
	sigs.k8s.io/cli-utils v0.22.1
	sigs.k8s.io/controller-runtime v0.6.4
)

replace (
	gopkg.in/square/go-jose.v2 => gopkg.in/square/go-jose.v2 v2.2.2
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.1-rc.0
	k8s.io/apiserver => k8s.io/apiserver v0.19.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.0
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.0
	k8s.io/code-generator => k8s.io/code-generator v0.19.1-rc.0
	k8s.io/component-base => k8s.io/component-base v0.19.0
	k8s.io/controller-manager => k8s.io/controller-manager v0.19.1-rc.0
	k8s.io/cri-api => k8s.io/cri-api v0.19.1-rc.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.0
	k8s.io/kubectl => k8s.io/kubectl v0.19.0
	k8s.io/kubelet => k8s.io/kubelet v0.19.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.0
	k8s.io/metrics => k8s.io/metrics v0.19.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.0
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.19.0
	k8s.io/sample-controller => k8s.io/sample-controller v0.19.0
)
