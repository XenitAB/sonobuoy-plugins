package fluxv2

import (
	"context"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/test/e2e/framework"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = ginkgo.Describe("fluxv2", func() {
	f := framework.NewDefaultFramework("fluxv2")

	var k8sClient client.Client

	ginkgo.BeforeEach(func() {
		err := sourcev1.AddToScheme(scheme.Scheme)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = kustomizev1.AddToScheme(scheme.Scheme)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		k8sClient, err = client.New(f.ClientConfig(), client.Options{Scheme: scheme.Scheme})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("should should successfully bootstrap", func() {
		ginkgo.By("checking the flux-system git source")
		key := types.NamespacedName{
			Name:      "flux-system",
			Namespace: "flux-system",
		}
		repo := &sourcev1.GitRepository{}
		err := k8sClient.Get(context.TODO(), key, repo)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(repo)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		res, err := status.Compute(&unstructured.Unstructured{Object: obj})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(res.Status.String()).To(gomega.Equal("Current"))

		ginkgo.By("checking the flux-system kustomization")
		key = types.NamespacedName{
			Name:      "flux-system",
			Namespace: "flux-system",
		}
		kust := &kustomizev1.Kustomization{}
		err = k8sClient.Get(context.TODO(), key, kust)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		obj, err = runtime.DefaultUnstructuredConverter.ToUnstructured(kust)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		res, err = status.Compute(&unstructured.Unstructured{Object: obj})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(res.Status.String()).To(gomega.Equal("Current"))
	})
})
