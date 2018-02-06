package kubemon

import (
	"time"

	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"

	// "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	_ "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func BuildAllInvolvedPodSelector() fields.Selector {
	return fields.ParseSelectorOrDie("involvedObject.kind=Pod")
}

func BuildInvolvedPodSelector(podName string) fields.Selector {
	return fields.AndSelectors(
		fields.ParseSelectorOrDie("involvedObject.kind=Pod"),
		fields.ParseSelectorOrDie("involvedObject.name="+podName))
}

type Watcher struct {
	clientset *kubernetes.Clientset
}

func WatchNodes(clientset *kubernetes.Clientset, sel fields.Selector, funcs cache.ResourceEventHandlerFuncs) (cache.Store, cache.Controller) {
	// resource nodes should give All namespaces
	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "nodes", core_v1.NamespaceAll, sel)
	return cache.NewInformer(
		watchlist,
		&core_v1.Node{},
		time.Minute*3,
		funcs)
}

func WatchEvents(clientset *kubernetes.Clientset, namespace string, sel fields.Selector, funcs cache.ResourceEventHandlerFuncs) (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "events", namespace, sel)
	return cache.NewInformer(
		watchlist,
		&core_v1.Event{},
		time.Minute*3,
		funcs)
}

// namespace can be "core_v1.NamespaceDefault"
func WatchJobs(clientset *kubernetes.Clientset, namespace string, sel fields.Selector, funcs cache.ResourceEventHandlerFuncs) (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(clientset.BatchV1().RESTClient(), "jobs", namespace, sel)
	return cache.NewInformer(
		watchlist,
		&batch_v1.Job{},
		time.Minute*3,
		funcs)

}

func WatchPods(clientset *kubernetes.Clientset, namespace string, sel fields.Selector, funcs cache.ResourceEventHandlerFuncs) (cache.Store, cache.Controller) {
	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", namespace, sel)
	return cache.NewInformer(
		watchlist,
		&core_v1.Pod{},
		time.Minute*3,
		funcs)
}

func FindPod(clientset *kubernetes.Clientset, namespace string, podId string) (*core_v1.Pod, error) {
	return clientset.CoreV1().Pods(namespace).Get(podId, meta_v1.GetOptions{})
}

func GetPods(clientset *kubernetes.Clientset, namespace string) (*core_v1.PodList, error) {
	return clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{})
}

func GetNodes(clientset *kubernetes.Clientset) (*core_v1.NodeList, error) {
	return clientset.CoreV1().Nodes().List(meta_v1.ListOptions{})
}
