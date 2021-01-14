package ingress

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/onsi/ginkgo"

	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	e2econfig "k8s.io/kubernetes/test/e2e/framework/config"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

var ingressConfig struct {
	Host string `usage:"hostname to use when creating ingress"`
}
var _ = e2econfig.AddOptions(&ingressConfig, "ingress")

var _ = ginkgo.Describe("Ingress TLS [Feature:Ingress]", func() {
	var ns string
	var c clientset.Interface
	f := framework.NewDefaultFramework("ingress")
	ginkgo.BeforeEach(func() {
		if len(ingressConfig.Host) == 0 {
			panic("ingress.host is required to be set")
		}

		c = f.ClientSet
		ns = f.Namespace.Name
	})

	ginkgo.It("should support tls ingress", func() {
		fmt.Println()

		ginkgo.By("creating a test pod")
		pod, err := c.CoreV1().Pods(ns).Create(context.TODO(), &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ingress-tls",
				Labels: map[string]string{
					"app": "ingress-tls",
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "server",
						Image: imageutils.GetE2EImage(imageutils.NginxNew),
					},
				},
			},
		}, metav1.CreateOptions{})
		framework.ExpectNoError(err)
		framework.Logf("Waiting for server to come up.")
		err = e2epod.WaitForPodRunningInNamespace(f.ClientSet, pod)
		framework.ExpectNoError(err, "Error occurred while waiting for pod status in namespace: Running.")

		ginkgo.By("creating a service for the pod")
		svc, err := c.CoreV1().Services(ns).Create(context.TODO(), &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ingress-tls",
			},
			Spec: v1.ServiceSpec{
				Ports: []v1.ServicePort{
					{
						Name:       "http",
						Port:       80,
						TargetPort: intstr.IntOrString{IntVal: 80},
					},
				},
				Selector: map[string]string{
					"app": "ingress-tls",
				},
			},
		}, metav1.CreateOptions{})
		framework.ExpectNoError(err)

		tlsIngressHost := fmt.Sprintf("ingress-tls.%s", ingressConfig.Host)
		ginkgo.By("creating a ingress with tls configuration")
		prefixPathType := extensionsv1beta1.PathTypePrefix
		_, err = c.ExtensionsV1beta1().Ingresses(ns).Create(context.TODO(), &extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ingress-tls",
				Labels: map[string]string{
					"app": "ingress-tls",
				},
				Annotations: map[string]string{
					"cert-manager.io/cluster-issuer":   "letsencrypt",
					"kubernetes.io/ingress.allow-http": "false",
				},
			},
			Spec: extensionsv1beta1.IngressSpec{
				Rules: []extensionsv1beta1.IngressRule{
					{
						Host: tlsIngressHost,
						IngressRuleValue: extensionsv1beta1.IngressRuleValue{
							HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
								Paths: []extensionsv1beta1.HTTPIngressPath{
									{
										PathType: &prefixPathType,
										Path:     "/",
										Backend: extensionsv1beta1.IngressBackend{
											ServiceName: svc.Name,
											ServicePort: intstr.FromString("http"),
										},
									},
								},
							},
						},
					},
				},
				TLS: []extensionsv1beta1.IngressTLS{
					{
						Hosts:      []string{tlsIngressHost},
						SecretName: "ingress-tls-cert",
					},
				},
			},
		}, metav1.CreateOptions{})
		framework.ExpectNoError(err)

		ginkgo.By("expecting https request to work")
		err = wait.Poll(10*time.Second, 15*time.Minute, func() (bool, error) {
			resp, err := http.Get(fmt.Sprintf("https://%s", tlsIngressHost))
			if err != nil {
				return false, nil
			}

			if resp.StatusCode != http.StatusOK {
				return false, nil
			}

			return true, nil
		})
		framework.ExpectNoError(err)

		ginkgo.By("expecting https redirect")
		err = wait.Poll(10*time.Second, 1*time.Minute, func() (bool, error) {
			resp, err := http.Get(fmt.Sprintf("http://%s", tlsIngressHost))
			if err != nil {
				return false, err
			}

			if resp.StatusCode != http.StatusOK {
				return false, nil
			}

			if resp.TLS == nil {
				return false, nil
			}

			return true, nil
		})
		framework.ExpectNoError(err)
	})
})
