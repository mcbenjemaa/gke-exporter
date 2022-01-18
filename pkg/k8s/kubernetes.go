package k8s

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func OutOfClusterClient() (*kubernetes.Clientset, error) {
	var kubeconfig string = filepath.Join(homedir.HomeDir(), ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Unable to initialise kube config %v", err)
	}

	// create the clientset
	return kubernetes.NewForConfig(config)
}

func GetClient() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Unable to initialise InClusterConfig %v", err)
	}
	// creates the clientset
	return kubernetes.NewForConfig(config)
}

func CountServices(ctx context.Context, clientset *kubernetes.Clientset) (int, error) {
	list, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("unable to list services %w", err)
	}

	return len(list.Items), nil
}

func CountPods(ctx context.Context, clientset *kubernetes.Clientset) (int, error) {
	list, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("unable to list pods %w", err)
	}

	return len(list.Items), nil
}
