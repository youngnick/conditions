package conditions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	v1 "github.com/youngnick/conditions/apis/youngnick.dev/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func GetConditions(configFlags *genericclioptions.ConfigFlags, objecttype string) ([]*v1.DuckConditions, error) {

	var returnvals []*v1.DuckConditions

	f := cmdutil.NewFactory(configFlags)

	// condConvert, err := NewUnstructuredConverter()
	// if err != nil {
	// 	return nil, err
	// }

	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return nil, err
	}

	client := dynamicClient.Resource(corev1.Resource(objecttype).WithVersion("v1"))

	result, err := client.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error pod %v", err)
	}

	for _, item := range result.Items {

		// obj, err := condConvert.FromUnstructured(item)
		// if err != nil {
		// 	return nil, err
		// }

		fmt.Printf("Object: %s/%s %s\n", item.GetNamespace(), item.GetName(), item.GetKind())

		raw, err := json.MarshalIndent(item.Object, "", "  ")
		if err != nil {
			return nil, err
		}

		// fmt.Println(string(raw))

		cond := &v1.DuckConditions{}

		err = json.Unmarshal(raw, cond)
		if err != nil {
			return nil, err
		}

		spew.Dump(cond)

		returnvals = append(returnvals, cond)

	}
	// config, err := configFlags.ToRESTConfig()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to read kubeconfig")
	// }

	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to create clientset")
	// }

	// _, resourcelists, err := clientset.DiscoveryClient.ServerGroupsAndResources()
	// if err != nil {
	// 	return nil, err
	// }

	// // for i, group := range groups {
	// // 	fmt.Printf("Group #%d\n", i)
	// // 	spew.Dump(group)
	// // }

	// for _, resourcelist := range resourcelists {

	// 	for _, resource := range resourcelist.APIResources {
	// 		if resource.Name == objecttype {
	// 			fmt.Printf("Resources: %s\n", resourcelist.GroupVersion)
	// 			spew.Dump(resource)
	// 		}
	// 	}
	// }
	// restClient := clientset.RESTClient()
	// config.ContentConfig.GroupVersion = &schema.GroupVersion{
	// 	Group:   "",
	// 	Version: "v1",
	// }
	// // config.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(info runtime.SerializerInfo)
	// spew.Dump(config)

	// // spew.Dump(restClient)
	// result := restClient.Get().Resource(objecttype).Do(context.Background())
	// // spew.Dump(result)

	// var conditionsslice []v1.DuckConditions

	// condition := &v1.DuckConditions{}

	// err = result.Into(condition)
	// if err != nil {
	// 	return nil, err
	// }
	// spew.Dump(condition)
	// // namespaces, err :=
	// CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to list namespaces")
	// }

	return returnvals, nil
}
