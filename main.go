package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Clone the repo and return the path
func cloneRepo() string {
	return "/foo"
}

// Store global k8s objects
func storeGlobals(client rest.Interface, path string) {

}

// Get a list of namespaces
func getNamespaces(client rest.Interface) []string {
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
	var client = clientset.Core().RESTClient()
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
