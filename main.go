package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

// Clone the repo and return the path
func cloneRepo() string {
	return "/foo"
}

func storePersistentVolumes(client corev1.CoreV1Interface, path string) {
	pvList, err := client.PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pv := range pvList.Items {
		fmt.Println("Found ", pv.Name)
	}
}

// Store global k8s objects
func storeGlobals(client corev1.CoreV1Interface, path string) {
	storePersistentVolumes(client, path)

}

// Get a list of namespaces
func getNamespaces(client corev1.CoreV1Interface) []string {
	return make([]string, 3)
}

// Store namespace objects
func storeNamespaces(namespace string, path string) {

}

// Make a git commit
func gitCommit(path string) {

}

// Push git commit
func gitPush(path string) {

}

func saveClusterState(clientset *kubernetes.Clientset) {
	var client = clientset.CoreV1()
	path := cloneRepo()
	storeGlobals(client, path)
	namespaces := getNamespaces(client)
	for _, namespace := range namespaces {
		storeNamespaces(namespace, path)
	}
	gitCommit(path)
	gitPush(path)
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// connect to k8s
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// Read k8k8state2git config here

	saveClusterState(clientset)

}
