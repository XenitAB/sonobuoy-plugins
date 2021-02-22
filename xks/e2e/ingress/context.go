package ingress

import (
	e2econfig "k8s.io/kubernetes/test/e2e/framework/config"
)

var ingressConfig struct {
	Host string `usage:"hostname to use when creating ingress"`
}
var _ = e2econfig.AddOptions(&ingressConfig, "ingress")
