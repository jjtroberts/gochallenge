package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	//apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var fYamlPath = flag.String("yamlpath", "", "specify a yaml resource")
var fImage = flag.String("image", "", "specify an image label")

func main() {
	flag.Parse()
	if *fYamlPath == "" {
		log.Fatal("must specify -yamlpath")
	}

	if *fImage == "" {
		log.Fatal("must specify -image")
	}

	if err := run(context.Background(), *fYamlPath, *fImage); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, yamlpath string, newImage string) error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Load the file into a buffer
	data, err := ioutil.ReadFile(yamlpath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a runtime.Decoder from the Codecs field within
	// k8s.io/client-go that's pre-loaded with the schemas for all
	// the standard Kubernetes resource types.
	decoder := scheme.Codecs.UniversalDeserializer()

	for _, resourceYAML := range strings.Split(string(data), "---") {

		// skip empty documents, `Decode` will fail on them
		if len(resourceYAML) == 0 {
			continue
		}

		// - obj is the API object (e.g., Deployment)
		// - groupVersionKind is a generic object that allows
		//   detecting the API type we are dealing with, for
		//   accurate type casting later.
		obj, groupVersionKind, err := decoder.Decode(
			[]byte(resourceYAML),
			nil,
			nil)
		if err != nil {
			log.Print(err)
			continue
		}

		// Figure out from `Kind` the resource type, and attempt
		// to cast appropriately.
		if groupVersionKind.Group == "apps" &&
			groupVersionKind.Version == "v1" &&
			groupVersionKind.Kind == "Deployment" {
			deployment := obj.(*appsv1.Deployment)

			depName := deployment.ObjectMeta.Name
			depNs := "default"

			if deployment.ObjectMeta.Namespace != "" {
				depNs = deployment.ObjectMeta.Namespace
			}

			log.Print("Deployment Name: " + depName)
			log.Print("Deployment NS: " + depNs)
			log.Print("Deployment Image: updating to " + newImage)

			deployment.Spec.Template.Spec.Containers[0].Image = newImage

			deploymentsClient := clientset.AppsV1().Deployments(depNs)

			_, getErr := deploymentsClient.Get(context.TODO(), depName, metav1.GetOptions{})

			// If not found
			if getErr != nil {
				fmt.Println("Creating deployment...")
				result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
				if err != nil {
					panic(err)
				}
				fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
			} else {
				fmt.Println("Updating deployment...")
				result, err := deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
				fmt.Printf("Updated deployment %q.\n", result.GetObjectMeta().GetName())
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return nil
}
