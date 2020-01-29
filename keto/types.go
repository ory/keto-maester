package keto

// ORYAccessControlPolicyJSON represents an OAuth2 client digestible by ORY Hydra
type ORYAccessControlPolicyJSON struct {
	ID          string                            `json:"id,omitempty"`
	Actions     []string                          `json:"actions,omitempty"`
	Conditions  map[string]map[string]interface{} `json:"conditions,omitempty"`
	Description string                            `json:"description"`
	Effect      string                            `json:"effect"`
	Resources   []string                          `json:"resources,omitempty"`
	Subjects    []string                          `json:"subjects,omitempty"`
}

// ORYAccessControlPolicyRoleJSON represents an OAuth2 client digestible by ORY Hydra
type ORYAccessControlPolicyRoleJSON struct {
	ID      string   `json:"id,omitempty"`
	Members []string `json:"members,omitempty"`
}
