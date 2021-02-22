package ingress

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	cmacmev1 "github.com/jetstack/cert-manager/pkg/apis/acme/v1"
	cmv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	imageutils "k8s.io/kubernetes/test/utils/image"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = ginkgo.Describe("Ingress TLS [Feature:Ingress]", func() {
	var ns string
	var clientSet clientset.Interface
	var k8sClient client.Client

	f := framework.NewDefaultFramework("ingress")
	ginkgo.BeforeEach(func() {
		if len(ingressConfig.Host) == 0 {
			panic("ingress.host is required to be set")
		}

		clientSet = f.ClientSet
		ns = f.Namespace.Name

		err := cmv1.AddToScheme(scheme.Scheme)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		k8sClient, err = client.New(f.ClientConfig(), client.Options{Scheme: scheme.Scheme})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("should support tls ingress", func() {
		ginkgo.By("creating a cluster issuer")
		clusterIssuer := &cmv1.ClusterIssuer{}
		nn := types.NamespacedName{
			Name: "letsencrypt",
		}
		err := k8sClient.Get(context.TODO(), nn, clusterIssuer)
		framework.ExpectNoError(err, "Could not get main clusterIssuer.")
		stagingClusterIssuer := &cmv1.ClusterIssuer{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-staging", nn.Name),
			},
			Spec: cmv1.IssuerSpec{
				IssuerConfig: cmv1.IssuerConfig{
					ACME: &cmacmev1.ACMEIssuer{
						Server:  "https://acme-staging-v02.api.letsencrypt.org/directory",
						Solvers: clusterIssuer.Spec.ACME.Solvers,
						PrivateKey: cmmetav1.SecretKeySelector{
							LocalObjectReference: cmmetav1.LocalObjectReference{
								Name: fmt.Sprintf("%s-private-key", nn.Name),
							},
						},
					},
				},
			},
		}
		err = k8sClient.Create(context.TODO(), stagingClusterIssuer)
		framework.ExpectNoError(err, "Could not create certmanager issuer.")
		defer k8sClient.Delete(context.TODO(), stagingClusterIssuer)

		name := fmt.Sprintf("%s-tls-test", ns)

		ginkgo.By("creating a test pod")
		pod, err := clientSet.CoreV1().Pods(ns).Create(context.TODO(), &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"app": name,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
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
		svc, err := clientSet.CoreV1().Services(ns).Create(context.TODO(), &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:       "http",
						Port:       80,
						TargetPort: intstr.IntOrString{IntVal: 80},
					},
				},
				Selector: map[string]string{
					"app": name,
				},
			},
		}, metav1.CreateOptions{})
		framework.ExpectNoError(err)

		ginkgo.By("creating a ingress with tls configuration")
		tlsIngressHost := fmt.Sprintf("%s.%s", name, ingressConfig.Host)
		prefixPathType := extensionsv1beta1.PathTypePrefix
		_, err = clientSet.ExtensionsV1beta1().Ingresses(ns).Create(context.TODO(), &extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"app": name,
				},
				Annotations: map[string]string{
					"cert-manager.io/cluster-issuer":   stagingClusterIssuer.Name,
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
						SecretName: fmt.Sprintf("%s-cert", name),
					},
				},
			},
		}, metav1.CreateOptions{})
		framework.ExpectNoError(err)

		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{Transport: tr}
		ginkgo.By("expecting https request to work")
		err = wait.Poll(20*time.Second, 15*time.Minute, func() (bool, error) {
			resp, err := client.Get(fmt.Sprintf("https://%s", tlsIngressHost))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
			gomega.Expect(resp.TLS).ShouldNot(gomega.BeNil())
			gomega.Expect(len(resp.TLS.VerifiedChains)).ShouldNot(gomega.Equal(0))
			gomega.Expect(resp.TLS.VerifiedChains[0][0].Subject.CommonName).Should(gomega.Equal(tlsIngressHost))
			gomega.Expect(resp.TLS.VerifiedChains[0][0].Issuer.CommonName).Should(gomega.Equal("Fake LE Intermediate X1"))
			return true, nil
		})
		framework.ExpectNoError(err)
		ginkgo.By("expecting https redirect")
		err = wait.Poll(20*time.Second, 1*time.Minute, func() (bool, error) {
			resp, err := client.Get(fmt.Sprintf("http://%s", tlsIngressHost))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
			gomega.Expect(resp.TLS).ShouldNot(gomega.BeNil())
			gomega.Expect(len(resp.TLS.VerifiedChains)).ShouldNot(gomega.Equal(0))
			gomega.Expect(resp.TLS.VerifiedChains[0][0].Subject.CommonName).Should(gomega.Equal(tlsIngressHost))
			gomega.Expect(resp.TLS.VerifiedChains[0][0].Issuer.CommonName).Should(gomega.Equal("Fake LE Intermediate X1"))
			return true, nil
		})
		framework.ExpectNoError(err)
	})
})
