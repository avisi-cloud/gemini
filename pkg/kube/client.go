package kube

import (
	"time"

	snapshotgroupv1 "github.com/fairwindsops/photon/pkg/types/snapshotgroup/v1"
	snapshotGroupClientset "github.com/fairwindsops/photon/pkg/types/snapshotgroup/v1/apis/clientset/versioned"
	"github.com/fairwindsops/photon/pkg/types/snapshotgroup/v1/apis/informers/externalversions"
	informers "github.com/fairwindsops/photon/pkg/types/snapshotgroup/v1/apis/informers/externalversions/snapshotgroup/v1"

	snapshotclient "github.com/kubernetes-csi/external-snapshotter/pkg/client/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Client provides access to k8s resources
type Client struct {
	K8s             kubernetes.Interface
	ClientSet       snapshotGroupClientset.Interface
	Informer        informers.SnapshotGroupInformer
	InformerFactory externalversions.SharedInformerFactory
	SnapshotClient  snapshotclient.Interface
}

var singleton *Client

// GetClient creates a new Client singleton
func GetClient() *Client {
	if singleton == nil {
		singleton = createClient()
	}
	return singleton
}

func createClient() *Client {
	kubeConf, configError := config.GetConfig()
	if configError != nil {
		panic(configError)
	}
	k8s, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		panic(err)
	}
	extClientSet, err := apiextensionsclient.NewForConfig(kubeConf)
	if err != nil {
		panic(err)
	}
	clientSet, err := snapshotGroupClientset.NewForConfig(kubeConf)
	if err != nil {
		panic(err)
	}
	snapshotClientSet, err := snapshotclient.NewForConfig(kubeConf)
	if err != nil {
		panic(err)
	}

	informerFactory := externalversions.NewSharedInformerFactory(clientSet, time.Second*30)
	informer := informerFactory.Snapshotgroup().V1().SnapshotGroups()

	if _, err = snapshotgroupv1.CreateCustomResourceDefinition("crd-ns", extClientSet); err != nil {
		panic(err)
	}
	return &Client{
		K8s:             k8s,
		ClientSet:       clientSet,
		Informer:        informer,
		InformerFactory: informerFactory,
		SnapshotClient:  snapshotClientSet,
	}
}
