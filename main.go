package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	itemv1 "github.com/keith-cullen/clientset/pkg/apis/itemresource/v1"
	myclientset "github.com/keith-cullen/clientset/pkg/client/clientset/versioned"
	myinformers "github.com/keith-cullen/clientset/pkg/client/informers/externalversions"
)

const (
	namespace = "itemns"
	timeout = 60 * time.Second
)

var (
	flags *flag.FlagSet
	logger klog.Logger
)

func usage() {
        writer := flags.Output()
        fmt.Fprintf(writer, "Usage: %s [OPTIONS]...\n", os.Args[0])
        flags.PrintDefaults()
}

func main() {
	flags = flag.NewFlagSet("", flag.ExitOnError)
	kubeconfigPath := flags.String("f", os.Getenv("KUBECONFIG"), "filename")
	flags.Usage = usage
	flags.Parse(os.Args[1:])  // ExitOnError so no need to check the return value
	logger = klog.LoggerWithName(klog.FromContext(context.Background()), "myclientset")
	logger.Info("main.main")
	kubernetesMasterEnv := os.Getenv("KUBERNETES_MASTER")
	kubernetesServiceHostEnv := os.Getenv("KUBERNETES_SERVICE_HOST")
	kubernetesServicePortEnv := os.Getenv("KUBERNETES_SERVICE_PORT")
	logger.Info("main.main", "kubeconfigPath", *kubeconfigPath)
	logger.Info("main.main", "KUBERNETES_MASTER", kubernetesMasterEnv)
	logger.Info("main.main", "KUBERNETES_SERVICE_HOST", kubernetesServiceHostEnv)
	logger.Info("main.main", "KUBERNETES_SERVICE_PORT", kubernetesServicePortEnv)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
	if err != nil {
		logger.Error(err, "Failed to create config")
		os.Exit(1)
	}

	// k8s clientset
	k8sClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err, "Failed to create k8s clientset")
		os.Exit(1)
	}
	pods, err := k8sClientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(err, "Failed to list pods")
		os.Exit(1)
	}
	logger.Info("main.main", "len(pods.Items)", len(pods.Items))

	// myclientset clientset
	myClientset, err := myclientset.NewForConfig(config)
	if err != nil {
		logger.Error(err, "Failed to create my clientset")
		os.Exit(1)
	}
	item1 := itemv1.Item{
		ObjectMeta: metav1.ObjectMeta{
			Name: "item1",
		},
		Spec: itemv1.ItemSpec{
			Detail: "xyz",
		},
	}
	_, err = myClientset.ItemresourceV1().Items(namespace).Create(context.TODO(), &item1, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Failed to create item")
		os.Exit(1)
	}
	items, err := myClientset.ItemresourceV1().Items(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(err, "Failed to list items")
		os.Exit(1)
	}
	logger.Info("main.main", "len(items.Items)", len(items.Items))
	for i, item := range items.Items {
		logger.Info("main.main", "item index", i, "Name", item.ObjectMeta.Name, "Detail", item.Spec.Detail)
	}

	// k8s informer
	k8sFactory := informers.NewSharedInformerFactory(k8sClientset, timeout)
	k8sPodInformer := k8sFactory.Core().V1().Pods()
	k8sCache := cache.ResourceEventHandlerFuncs{
		AddFunc:    k8sPodAdd,
		UpdateFunc: k8sPodUpdate,
		DeleteFunc: k8sPodDelete,
	}
	k8sPodInformer.Informer().AddEventHandler(k8sCache)
	var k8sStopCh chan struct{}
	k8sFactory.Start(k8sStopCh)
	if !cache.WaitForCacheSync(k8sStopCh, k8sPodInformer.Informer().HasSynced) {
		logger.Error(err, "Failed to sync k8s informer")
		os.Exit(1)
	}

	// myclientset informer
	myFactory := myinformers.NewSharedInformerFactory(myClientset, timeout)
	myItemInformer := myFactory.Itemresource().V1().Items()
	myCache := cache.ResourceEventHandlerFuncs{
		AddFunc:    myItemAdd,
		UpdateFunc: myItemUpdate,
		DeleteFunc: myItemDelete,
	}
	myItemInformer.Informer().AddEventHandler(myCache)
	var myStopCh chan struct{}
	myFactory.Start(myStopCh)
	if !cache.WaitForCacheSync(myStopCh, myItemInformer.Informer().HasSynced) {
		logger.Error(err, "Failed to sync my informer")
		os.Exit(1)
	}
}

func k8sPodAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	logger.Info("main.k8sPodAdd",
		"Namespace", pod.Namespace,
		"Name", pod.Name)
}

func k8sPodUpdate(oldObj, newObj interface{}) {
	oldPod := oldObj.(*v1.Pod)
	newPod := newObj.(*v1.Pod)
	logger.Info("main.k8sPodUpdate",
		"Old Namespace", oldPod.Namespace,
		"Old Name", oldPod.Name,
		"New Namespace", newPod.Namespace,
		"New Name", newPod.Name)
}

func k8sPodDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	logger.Info("main.k8sPodDelete",
		"Namespace", pod.Namespace,
		"Name", pod.Name)
}

func myItemAdd(obj interface{}) {
	item := obj.(*itemv1.Item)
	logger.Info("main.myItemAdd",
		"Namespace", item.Namespace,
		"Name", item.Name)
}

func myItemUpdate(oldObj, newObj interface{}) {
	oldItem := oldObj.(*itemv1.Item)
	newItem := newObj.(*itemv1.Item)
	logger.Info("main.myItemUpdate",
		"Old Namespace", oldItem.Namespace,
		"Old Name", oldItem.Name)
		"New Namespace", newItem.Namespace,
		"New Name", newItem.Name)
	logger.Info("main.myItemUpdate")
}

func myItemDelete(obj interface{}) {
	item := obj.(*itemv1.Item)
	logger.Info("main.myItemDelete",
		"Namespace", item.Namespace,
		"Name", item.Name)
}
