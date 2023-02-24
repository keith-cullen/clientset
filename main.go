package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	itemv1 "github.com/keith-cullen/clientset/pkg/apis/itemresource/v1"
	myclientset "github.com/keith-cullen/clientset/pkg/client/clientset/versioned"
)

const namespace = "itemns"

var flags *flag.FlagSet

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
	logger := klog.LoggerWithName(klog.FromContext(context.Background()), "myclientset")
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
	k8s_clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err, "Failed to create k8s clientset")
		os.Exit(1)
	}
	pods, err := k8s_clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(err, "Failed to list pods")
		os.Exit(1)
	}
	logger.Info("main.main", "len(pods.Items)", len(pods.Items))
	my_clientset, err := myclientset.NewForConfig(config)
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
	_, err = my_clientset.ItemresourceV1().Items(namespace).Create(context.TODO(), &item1, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Failed to create item")
		os.Exit(1)
	}
	items, err := my_clientset.ItemresourceV1().Items(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(err, "Failed to list items")
		os.Exit(1)
	}
	logger.Info("main.main", "len(items.Items)", len(items.Items))
	for i, item := range items.Items {
		logger.Info("main.main", "item index", i, "Name", item.ObjectMeta.Name, "Detail", item.Spec.Detail)
	}
}
