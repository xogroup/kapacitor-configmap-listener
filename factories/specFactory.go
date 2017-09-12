package factories

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api/v1"
)

// SpecFactory is an object to hold any api/v1 spec builders
type SpecFactory struct{}

// NewSpecFactory instantiates a new object of that type
func NewSpecFactory() *SpecFactory {
	return &SpecFactory{}
}

// Build outputs an api/v1 resource struct matching the input resource name
func (*SpecFactory) Build(resource string) runtime.Object {

	switch resource {
	case "services":
		return &v1.Service{}
	case "configmaps":
		return &v1.ConfigMap{}
	}

	panic(fmt.Errorf("no resource mapped for %s", resource))
}
