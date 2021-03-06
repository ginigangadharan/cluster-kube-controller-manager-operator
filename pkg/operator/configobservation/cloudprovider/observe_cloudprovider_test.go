package cloudprovider

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/openshift/library-go/pkg/operator/events"

	"github.com/openshift/cluster-kube-controller-manager-operator/pkg/operator/configobservation"
)

func TestObserveCloudProviderNames(t *testing.T) {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	if err := indexer.Add(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-config-v1",
			Namespace: "kube-system",
		},
		Data: map[string]string{
			"install-config": "platform:\n  aws: {}\n",
		},
	}); err != nil {
		t.Fatal(err.Error())
	}
	listers := configobservation.Listers{
		ConfigmapLister: corelistersv1.NewConfigMapLister(indexer),
	}
	result, errs := ObserveCloudProviderNames(listers, events.NewInMemoryRecorder("cloud"), map[string]interface{}{})
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	cloudProvider, _, err := unstructured.NestedSlice(result, "extendedArguments", "cloud-provider")
	if err != nil {
		t.Fatal(err)
	}
	if e, a := 1, len(cloudProvider); e != a {
		t.Fatalf("expected len(cloudProvider) == %d, got %d", e, a)
	}
	if e, a := "aws", cloudProvider[0]; e != a {
		t.Errorf("expected cloud-provider=%s, got %s", e, a)
	}
}
