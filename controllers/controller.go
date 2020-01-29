package controllers

import (
	ketov1alpha1 "github.com/ory/keto-maester/api/v1alpha1"
	"github.com/ory/keto-maester/keto"
)

const (
	FinalizerName = "finalizer.ory.keto.sh"
)

type KetoConfiger interface {
	GetKeto() ketov1alpha1.Keto
}

type KetoClientMakerFunc func(KetoConfiger) (KetoClientInterface, error)

type clientMapKey struct {
	url      string
	port     int
	endpoint string
}

type KetoClientInterface interface {
	GetORYAccessControlPolicy(flavor, id string) (*keto.ORYAccessControlPolicyJSON, bool, error)
	ListORYAccessControlPolicy(flavor string) ([]*keto.ORYAccessControlPolicyJSON, error)
	PutORYAccessControlPolicy(flavor string, o *keto.ORYAccessControlPolicyJSON) (*keto.ORYAccessControlPolicyJSON, error)
	DeleteORYAccessControlPolicy(flavor, id string) error
	GetORYAccessControlPolicyRole(flavor, id string) (*keto.ORYAccessControlPolicyRoleJSON, bool, error)
	ListORYAccessControlPolicyRole(flavor string) ([]*keto.ORYAccessControlPolicyRoleJSON, error)
	PutORYAccessControlPolicyRole(flavor string, o *keto.ORYAccessControlPolicyRoleJSON) (*keto.ORYAccessControlPolicyRoleJSON, error)
	DeleteORYAccessControlPolicyRole(flavor, id string) error
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
