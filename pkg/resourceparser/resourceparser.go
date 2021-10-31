// package resourceparser parses object names passed on the command line and
// dereferences them to a GVR and other info.
package resourceparser

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type Resource struct {
	Name  string
	Group string
	GVR   schema.GroupVersionResource
}

func Translate(f cmdutil.Factory, res string) (*schema.GroupVersionResource, error) {

	// TODO: Refactor building the lookup table to be a separate function, and have Translate
	// just lookup in the lookup table.
	discoveryclient, err := f.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	lists, err := discoveryclient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	lookupTable := make(resourceLookup)

	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			continue
		}

		for _, resource := range list.APIResources {
			rn := resource.Name

			if len(resource.Verbs) == 0 {
				continue
			}

			gvr := gv.WithResource(rn)

			lookupTable[resource.Name] = gvr

			lookupTable[strings.ToLower(resource.Kind)] = gvr
			if resource.SingularName != "" {
				lookupTable[resource.SingularName] = gvr
			}

			for _, short := range resource.ShortNames {
				lookupTable[short] = gvr
			}

		}
	}

	parsedGVR, ok := lookupTable[res]
	if !ok {
		return nil, fmt.Errorf("Couldn't find resources called %s", res)
	}

	return &parsedGVR, nil
}

// groupResource contains the APIGroup and APIResource
type groupResource struct {
	APIGroup        string
	APIGroupVersion string
	APIResource     metav1.APIResource
}

type resourceLookup map[string]schema.GroupVersionResource
