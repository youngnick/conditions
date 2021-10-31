package conditions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	v1 "github.com/youngnick/conditions/apis/youngnick.dev/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func GetConditions(f cmdutil.Factory, objecttype *schema.GroupVersionResource) ([]ConditionsObj, error) {

	var returnvals []ConditionsObj

	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return nil, err
	}

	client := dynamicClient.Resource(*objecttype)

	result, err := client.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error pod %v", err)
	}

	for _, item := range result.Items {

		var condObj ConditionsObj

		condObj.Name = item.GetName()
		condObj.Namespace = item.GetNamespace()
		condObj.Group = item.GetObjectKind().GroupVersionKind().Group
		condObj.Resource = objecttype.Resource
		condObj.APIVersion = item.GetAPIVersion()

		fmt.Printf("Object: %s/%s %s\n", item.GetNamespace(), item.GetName(), item.GetKind())

		raw, err := json.MarshalIndent(item.Object, "", "  ")
		if err != nil {
			return nil, err
		}

		cond := &v1.DuckConditions{}

		err = json.Unmarshal(raw, cond)
		if err != nil {
			return nil, err
		}

		condObj.Conditions = cond.Status.Conditions

		// spew.Dump(condObj)
		returnvals = append(returnvals, condObj)

	}
	return returnvals, nil
}

type ConditionsObj struct {
	Name       string
	Namespace  string
	Resource   string
	APIVersion string
	Group      string
	Conditions []metav1.Condition
}

func (co *ConditionsObj) Print(out io.Writer) {

	var outputName string

	outputName = co.Resource + "." + co.Group + "/" + co.Name

	fmt.Fprintf(out, "\n-n %s %s\n\nConditions:\n", co.Namespace, outputName)

	fmt.Fprintf(out, "  Type\tStatus\tLastTransitionTime\tReason\tMessage\n")
	fmt.Fprintf(out, "  ----\t------\t------------------\t------\t-------\n")

	for _, c := range co.Conditions {
		fmt.Fprintf(out, "  %v \t%v \t%s \t%v \t%v\n",
			c.Type,
			c.Status,
			c.LastTransitionTime.Time.Format(time.RFC1123Z),
			c.Reason,
			c.Message)
	}
}
