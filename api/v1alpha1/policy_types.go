/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/ory/keto-maester/keto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatusCode string

const (
	StatusUpsertPolicyFailed StatusCode = "POLICY_UPSERT_FAILED"
	StatusInvalidKetoAddress StatusCode = "INVALID_KETO_ADDRESS"
)

// Keto defines the desired keto instance to use for ORYAccessControlPolicy
type Keto struct {
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=(^$|^https?://.*)
	//
	// URL is the URL for the keto instance on
	// which to set up the client. This value will override the value
	// provided to `--keto-url`
	URL string `json:"url,omitempty"`

	// +kubebuilder:validation:Maximum=65535
	//
	// Port is the port for the keto instance on
	// which to set up the client. This value will override the value
	// provided to `--keto-port`
	Port int `json:"port,omitempty"`

	// +kubebuilder:validation:Pattern=(^$|^/.*)
	//
	// Endpoint is the endpoint for the keto instance on which
	// to set up the client. This value will override the value
	// provided to `--endpoint` (defaults to `"/engines"` in the
	// application)
	Endpoint string `json:"endpoint,omitempty"`
}

// GetKeto to satisfy the controllers.KetoConfiger interface
func (k Keto) GetKeto() Keto {
	return k
}

// Condition defines a condition
type Condition struct {
	// +kubebuilder:validation:Enum=CIDRCondition;StringEqualCondition;StringMatchCondition;EqualsSubjectCondition;StringPairsEqualCondition
	//
	// Type is the type of the condition
	Type string `json:"type"`

	// Options is the options for the condition
	Options map[string]string `json:"options,omitempty"`
}

// ORYAccessControlPolicySpec defines the desired state of OAuth2Client
type ORYAccessControlPolicySpec struct {
	// +kubebuilder:validation:Enum=exact;glob;regex
	//
	// Flavor is the flavor for the policy
	Flavor string `json:"flavor"`

	// +kubebuilder:validation:MinLength=1
	//
	// ID is the policy id (name). If this is not provided
	// `metadata.name` is used
	ID string `json:"id,omitempty"`

	// +kubebuilder:validation:MinItems=1
	//
	// Actions is an array of actions to which the policy applies.
	Actions []URN `json:"actions"`

	// Conditions is a conditions object to apply to the policy.
	Conditions map[string]Condition `json:"conditions,omitempty"`

	// Description is the policy description
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Enum=allow;deny
	//
	// Effect is the effect of the policy
	Effect string `json:"effect"`

	// +kubebuilder:validation:MinItems=1
	//
	// Resources is an array of resources to which the policy applies.
	Resources []URN `json:"resources"`

	// +kubebuilder:validation:MinItems=1
	//
	// Subjects is an array of subjects to which the policy applies.
	Subjects []Subject `json:"subjects"`

	// Keto is the optional configuration to use for managing
	// this policy
	Keto Keto `json:"keto,omitempty"`
}

// +kubebuilder:validation:MinLength=1
//
// URN is a resource or action URN
type URN string

// +kubebuilder:validation:MinLength=1
//
// Subject is a subject name
type Subject string

// ORYAccessControlPolicyStatus defines the observed state of OAuth2Client
type ORYAccessControlPolicyStatus struct {
	// ObservedGeneration represents the most recent generation observed by the daemon set controller.
	ObservedGeneration  int64               `json:"observedGeneration,omitempty"`
	ReconciliationError ReconciliationError `json:"reconciliationError,omitempty"`
}

// ReconciliationError represents an error that occurred during the reconciliation process
type ReconciliationError struct {
	// Code is the status code of the reconciliation error
	Code StatusCode `json:"statusCode,omitempty"`
	// Description is the description of the reconciliation error
	Description string `json:"description,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ORYAccessControlPolicy is the Schema for the ORYAccessControlPolicy API
type ORYAccessControlPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ORYAccessControlPolicySpec   `json:"spec,omitempty"`
	Status ORYAccessControlPolicyStatus `json:"status,omitempty"`
}

// GetKeto returns the keto config on the spec
func (o ORYAccessControlPolicy) GetKeto() Keto {
	return o.Spec.Keto
}

// +kubebuilder:object:root=true

// ORYAccessControlPolicyList contains a list of ORYAccessControlPolicy
type ORYAccessControlPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ORYAccessControlPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ORYAccessControlPolicy{}, &ORYAccessControlPolicyList{})
}

// ToORYAccessControlPolicyJSON converts an ORYAccessControlPolicy into a ORYAccessControlPolicyJSON object that represents an OAuth2 client digestible by ORY Hydra
func (c *ORYAccessControlPolicy) ToORYAccessControlPolicyJSON() *keto.ORYAccessControlPolicyJSON {
	id := c.ObjectMeta.Name
	if c.Spec.ID != "" {
		id = c.Spec.ID
	}

	return &keto.ORYAccessControlPolicyJSON{
		ID:          id,
		Actions:     urnToStringSlice(c.Spec.Actions),
		Conditions:  conditionsToMap(c.Spec.Conditions),
		Description: c.Spec.Description,
		Effect:      c.Spec.Effect,
		Resources:   urnToStringSlice(c.Spec.Resources),
		Subjects:    subjectsToStringSlice(c.Spec.Subjects),
	}
}

func urnToStringSlice(rt []URN) []string {
	var output = make([]string, len(rt))
	for i, elem := range rt {
		output[i] = string(elem)
	}
	return output
}

func subjectsToStringSlice(gt []Subject) []string {
	var output = make([]string, len(gt))
	for i, elem := range gt {
		output[i] = string(elem)
	}
	return output
}

func conditionsToMap(cs map[string]Condition) map[string]map[string]interface{} {
	conditions := map[string]map[string]interface{}{}

	for n, c := range cs {
		conditions[n] = map[string]interface{}{
			"type":    c.Type,
			"options": c.Options,
		}
	}

	return conditions
}
