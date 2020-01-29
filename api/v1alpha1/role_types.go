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

const (
	StatusUpsertRoleFailed StatusCode = "ROLE_UPSERT_FAILED"
)

// ORYAccessControlPolicyRoleSpec defines the desired state of OAuth2Client
type ORYAccessControlPolicyRoleSpec struct {
	// +kubebuilder:validation:Enum=exact;glob;regex
	//
	// Flavor is the flavor for the role
	Flavor string `json:"flavor"`

	// +kubebuilder:validation:MinLength=1
	//
	// ID is the role id (name). If this is not provided
	// `metadata.name` is used
	ID string `json:"id,omitempty"`

	// +kubebuilder:validation:MinItems=1
	//
	// Members is an array of members belonging to this role.
	Members []Subject `json:"members"`

	// Keto is the optional configuration to use for managing
	// this role
	Keto Keto `json:"keto,omitempty"`
}

// ORYAccessControlPolicyRoleStatus defines the observed state of OAuth2Client
type ORYAccessControlPolicyRoleStatus struct {
	// ObservedGeneration represents the most recent generation observed by the daemon set controller.
	ObservedGeneration  int64               `json:"observedGeneration,omitempty"`
	ReconciliationError ReconciliationError `json:"reconciliationError,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ORYAccessControlPolicyRole is the Schema for the ORYAccessControlPolicyRole API
type ORYAccessControlPolicyRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ORYAccessControlPolicyRoleSpec   `json:"spec,omitempty"`
	Status ORYAccessControlPolicyRoleStatus `json:"status,omitempty"`
}

// GetKeto returns the keto config on the spec
func (o ORYAccessControlPolicyRole) GetKeto() Keto {
	return o.Spec.Keto
}

// +kubebuilder:object:root=true

// ORYAccessControlPolicyRoleList contains a list of ORYAccessControlPolicyRole
type ORYAccessControlPolicyRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ORYAccessControlPolicyRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ORYAccessControlPolicyRole{}, &ORYAccessControlPolicyRoleList{})
}

// ToORYAccessControlPolicyRoleJSON converts an ORYAccessControlPolicyRole into a ORYAccessControlPolicyRoleJSON object that represents an OAuth2 client digestible by ORY Hydra
func (c *ORYAccessControlPolicyRole) ToORYAccessControlPolicyRoleJSON() *keto.ORYAccessControlPolicyRoleJSON {
	id := c.ObjectMeta.Name
	if c.Spec.ID != "" {
		id = c.Spec.ID
	}
	return &keto.ORYAccessControlPolicyRoleJSON{
		ID:      id,
		Members: subjectsToStringSlice(c.Spec.Members),
	}
}
