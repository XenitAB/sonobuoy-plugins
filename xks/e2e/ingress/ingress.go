package ingress

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/onsi/ginkgo"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

var _ = ginkgo.Describe("Ingress TLS", func() {
	var ns string
	var c clientset.Interface

	f := framework.NewDefaultFramework("ingress")

	ingressHost := ""

	ginkgo.BeforeEach(func() {
		c = f.ClientSet
		ns = f.Namespace.Name
	})

	ginkgo.It("should support tls ingress", func() {
		ginkgo.By("creating a test pod")
		pod, err := c.CoreV1().Pods(ns).Create(context.TODO(), &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "ingress-tls-",
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
				GenerateName: "ingress-tls-",
			},
			Spec: v1.ServiceSpec{
				Ports: []v1.ServicePort{
					{
						Name:       "http",
						Port:       80,
						TargetPort: intstr.IntOrString{IntVal: 8080},
					},
				},
				Selector: map[string]string{
					"app": "ingress-tls",
				},
			},
		}, metav1.CreateOptions{})
		framework.ExpectNoError(err)

		ginkgo.By("creating a ingress with tls configuration")
		prefixPathType := networkingv1.PathTypePrefix
		_, err = c.NetworkingV1().Ingresses(ns).Create(context.TODO(), &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "ingress-tls-",
				Labels: map[string]string{
					"app": "ingress-tls",
				},
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: ingressHost,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										PathType: &prefixPathType,
										Path:     "/",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Port: networkingv1.ServiceBackendPort{
													Name: "http",
												},
												Name: svc.Name,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, metav1.CreateOptions{})
		framework.ExpectNoError(err)

		ginkgo.By("expecting https request to work")
		err = wait.Poll(10*time.Second, 1*time.Minute, func() (bool, error) {
			resp, err := http.Get(fmt.Sprintf("https://%s", ingressHost))
			if err != nil {
				return false, err
			}

			if resp.StatusCode != http.StatusOK {
				return false, nil
			}

			return true, nil
		})
		framework.ExpectNoError(err)

		ginkgo.By("expecting https redirect to work")
		err = wait.Poll(10*time.Second, 1*time.Minute, func() (bool, error) {
			resp, err := http.Get(fmt.Sprintf("http://%s", ingressHost))
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
