package clusterhealth

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apipod "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/test/e2e/framework"
	e2enode "k8s.io/kubernetes/test/e2e/framework/node"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
)

var _ = ginkgo.Describe("Cluster health check [Feature:ClusterHealth]", func() {
	var clientSet clientset.Interface

	f := framework.NewDefaultFramework("cluster-health")
	ginkgo.BeforeEach(func() {
		clientSet = f.ClientSet
	})

	ginkgo.It("all pods in kube-system should be running [Troubleshoot]", func() {
		ns := "kube-system"
		podList, err := clientSet.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		framework.ExpectNoError(err, "could not list pods")
		podNames := []string{}
		for _, pod := range podList.Items {
			podNames = append(podNames, pod.Name)
		}
		e2epod.CheckPodsRunningReady(clientSet, ns, podNames, 5*time.Second)
	})

	ginkgo.It("checking all nodes are ready [Troubleshoot]", func() {
		nodeList, err := clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		framework.ExpectNoError(err, "could not list nodes")
		e2enode.CheckReady(clientSet, len(nodeList.Items), 5*time.Second)
	})

	ginkgo.It("checks unschedulable pods [Troubleshoot]", func() {
		podList, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		framework.ExpectNoError(err, "could not list pods")

		err = nil
		for _, pod := range podList.Items {
			index, condition := apipod.GetPodCondition(&pod.Status, corev1.PodScheduled)
			if index == -1 {
				continue
			}
			if condition.Reason == "Unschedulable" {
				err = multierr.Append(err, fmt.Errorf("pod %q could not be scheduled: %q", pod.Name, condition.Message))
			}
		}
		framework.ExpectNoError(err)
	})

	ginkgo.It("checks image pull errors [Troubleshoot]", func() {
		podList, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		framework.ExpectNoError(err, "could not list pods")

		err = nil
		for _, pod := range podList.Items {
      for _, container := range pod.Spec.Containers {
        status, present := apipod.GetContainerStatus(pod.Status.ContainerStatuses, container.Name)
        if !present {
          continue
        }

        if status.State.Waiting.Reason == "ImagePullBackOff" {
						
        }
      }
		}
		framework.ExpectNoError(err)
	})

	ginkgo.It("checks container restarts [Troubleshoot]", func() {

	})

	ginkgo.It("looks at pod distribution over nodes [Troubleshoot]", func() {

	})

	ginkgo.It("calculates average memory/cpu utilization [Troubleshoot]", func() {

	})

	/*ginkgo.It("memory request and limit set [Troubleshoot]", func() {
		podList, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		framework.ExpectNoError(err, "could not list pods")

		err = nil
		for _, pod := range podList.Items {
			if pod.Namespace == "kube-system" {
				continue
			}

			for _, container := range pod.Spec.Containers {
				if container.Resources.Requests.Memory().IsZero() || container.Resources.Limits.Memory().IsZero() {
					err = multierr.Append(err, fmt.Errorf("pod %q does not set memory request or limit, may cause uncontrolled resource consumption", pod.Name))
				}
			}
		}
		framework.ExpectNoError(err)
	})*/

