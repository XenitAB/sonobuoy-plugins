package highavailability

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
)

var _ = ginkgo.Describe("High availibiltiy configuration [Feature:HighAvailability]", func() {
	var clientSet clientset.Interface

	f := framework.NewDefaultFramework("high-availibiltiy")
	ginkgo.BeforeEach(func() {
		clientSet = f.ClientSet
	})

	ginkgo.It("deployment replicas should run on multiple nodes [Troubleshoot]", func() {
		depList, err := clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		framework.ExpectNoError(err, "could not list deployments")

		result := []string{}
		for _, dep := range depList.Items {
			if dep.Status.Replicas <= 1 {
				framework.Logf("skipping deployment %q as it only has a single replica", dep.Name)
				continue
			}

			framework.Logf("checking deployment %q replicas", dep.Name)
			labelMap, err := metav1.LabelSelectorAsMap(dep.Spec.Selector)
			framework.ExpectNoError(err, "could not create label selector map")
			opt := metav1.ListOptions{LabelSelector: labels.Set(labelMap).String()}
			podList, err := clientSet.CoreV1().Pods(dep.Namespace).List(context.TODO(), opt)
			framework.ExpectNoError(err, "could not list pods for deployment")
			if !runningOnMultipleNodes(podList.Items) {
				result = append(result, fmt.Sprintf("deployment %q with multiple replicas is only running on node %q", dep.Name, podList.Items[0].Spec.NodeName))
			}
		}
		framework.ExpectEmpty(result, "expected error list to be empty")
	})
})

// TODO (phillebaba) not important but can be optimized
// to avoid iterations
func runningOnMultipleNodes(pods []corev1.Pod) bool {
	for _, podOuter := range pods {
		for _, podInner := range pods {
			if podOuter.Spec.NodeName != podInner.Spec.NodeName {
				return true
			}
		}
	}

	return false
}
