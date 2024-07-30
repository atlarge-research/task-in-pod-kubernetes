package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string
var crdGVK = schema.GroupVersionKind{
	Group:   "example.com",
	Version: "v1",
	Kind:    "MyPod",
}

func main() {
	fmt.Printf("------INSIDE MAIN")
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	// Set up Kubernetes client
	config, err := getClientConfig(*kubeconfig)
	if err != nil {
		fmt.Printf("Error getting client config: %v\n", err)
		return
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes clientset: %v\n", err)
		return
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating dynamic client: %v\n", err)
		return
	}

	// Watch for changes to MyPod CR
	go watchMyPodChanges(dynamicClient, clientset)

	// Watch for changes to controller.yaml and apply changes
	go watchControllerChanges(clientset)

	// Wait forever
	select {}
}

func getClientConfig(kubeconfig string) (*rest.Config, error) {
	var config *rest.Config
	var err error
	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func watchMyPodChanges(client dynamic.Interface, clientset *kubernetes.Clientset) {
	fmt.Printf("+++++++++++ WATCHING FOR CHANGES")
	ctx := context.Background()
	for {
		watcher, err := client.Resource(schema.GroupVersionResource{
			Group:    crdGVK.Group,
			Version:  crdGVK.Version,
			Resource: "mypods",
		}).Watch(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error watching MyPod changes: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		for event := range watcher.ResultChan() {

			obj, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				fmt.Println("Error converting to unstructured:", err)
				continue
			}
			// Extract name from MyPod
			name := obj.GetName()

			//fmt.Printf("--------CREATE POD IS RUNNING----------- ", name)

			// Check the event type
			switch event.Type {
			case watch.Added:
				// Create a new Pod object
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name, // Set the name of the Pod
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						// Add Pod specifications as needed
						Containers: []corev1.Container{
							{
								Name:  "mypod-container",
								Image: "busybox", // Example image, change as needed
							},
						},
					},
				}

				// Apply the new Pod object
				_, err = clientset.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
				if err != nil {
					fmt.Printf("Error creating Pod: %v\n", err)
					continue
				}
			case watch.Deleted:
				// Delete the Pod with the given name
				err := clientset.CoreV1().Pods("default").Delete(ctx, name, metav1.DeleteOptions{})
				if err != nil {
					fmt.Printf("Error deleting Pod: %v\n", err)
					continue
				}
			}
		}
	}
}

func watchControllerChanges(clientset *kubernetes.Clientset) {
	fmt.Printf("#####WATCH CONTROLLER")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating file watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Add controller.yaml to watcher
	controllerFile := "controller.yaml"
	err = watcher.Add(controllerFile)
	if err != nil {
		fmt.Printf("Error adding file to watcher: %v\n", err)
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("Changes detected in controller.yaml. Applying changes...")
				// Apply kubectl apply command
				cmd := exec.Command("kubectl", "apply", "-f", controllerFile)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Printf("Error applying kubectl command: %v\n", err)
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("Error in file watcher: %v\n", err)
		}
	}
}
