package kubemon

import (
	"bitbucket.org/linkernetworks/aurora/src/kubeconfig"
	"bitbucket.org/linkernetworks/aurora/src/logger"

	"github.com/stretchr/testify/assert"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"os"
	"testing"
)

func TestGetPods(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_K8S"); !defined {
		t.SkipNow()
		return
	}

	config, err := kubeconfig.Load("", "")
	assert.NoError(t, err)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	assert.NoError(t, err)

	pods, err := GetPods(clientset, "testing")
	assert.NoError(t, err)
	assert.NotEmpty(t, pods)
	t.Log(pods)
}

func TestWatchEvents(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_K8S"); !defined {
		t.Skip("Skip kubernetes tests")
		return
	}

	config, err := kubeconfig.Load("", "")
	assert.NoError(t, err)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	assert.NoError(t, err)

	var stream = make(chan *core_v1.Event)

	// http://127.0.0.1:8001/api/v1/events?selector=involvedObject.name=job-5a795233ab63a766626a5e32-0-run-0-k84p5
	_, controller := WatchEvents(clientset, "default", BuildAllInvolvedPodSelector(),
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(a interface{}) {
				if event, ok := a.(*core_v1.Event); ok {
					stream <- event
				}
			},
			UpdateFunc: func(a, b interface{}) {
				if event, ok := a.(*core_v1.Event); ok {
					stream <- event
				}
			},
		})

	stop := make(chan struct{})

	go controller.Run(stop)

STREAM:
	for {
		select {
		case event := <-stream:
			logger.Infof("event: %+v", event)
			break STREAM
		}
	}

	var e struct{}
	stop <- e
}

func TestGetNodes(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_K8S"); !defined {
		t.SkipNow()
		return
	}

	config, err := kubeconfig.Load("", "")
	assert.NoError(t, err)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	assert.NoError(t, err)

	nodes, err := GetNodes(clientset)
	assert.NoError(t, err)
	assert.NotEmpty(t, nodes)
	t.Log(nodes)
}
