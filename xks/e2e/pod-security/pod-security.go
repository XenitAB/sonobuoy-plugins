package podsecurity

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

var _ = ginkgo.Describe("Pod container security [Feature:PodSecurity]", func() {
	var ns string
	var clientSet clientset.Interface

	f := framework.NewDefaultFramework("pod-security")
	ginkgo.BeforeEach(func() {
		clientSet = f.ClientSet
		ns = f.Namespace.Name
	})

	ginkgo.It("should not allow privileged containers", func() {
		name := fmt.Sprintf("%s", ns)
		pod := securityPod(name)
		pod.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{Privileged: boolPointer(true)}
		_, err := clientSet.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
		framework.ExpectError(err, "creating pods that ave privileged should not be allowed")
	})

	ginkgo.It("should not allow privlege escalation", func() {
		name := fmt.Sprintf("%s", ns)
		pod := securityPod(name)
		pod.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{AllowPrivilegeEscalation: boolPointer(true)}
		_, err := clientSet.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
		framework.ExpectError(err)
	})

	ginkgo.It("should not allow host path mounting", func() {
		name := fmt.Sprintf("%s", ns)
		pod := securityPod(name)
		pod.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
			{
				MountPath: "/tmp/root",
				Name:      "root",
			},
		}
		pod.Spec.Volumes = []corev1.Volume{
			{
				Name: "root",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/root",
					},
				},
			},
		}
		_, err := clientSet.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
		framework.ExpectError(err)
	})
})

func boolPointer(value bool) *bool {
	return &value
}

func securityPod(name string) *corev1.Pod {
	return &corev1.Pod{
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
	}
}
