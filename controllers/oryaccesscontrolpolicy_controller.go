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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	ketov1alpha1 "github.com/ory/keto-maester/api/v1alpha1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ORYAccessControlPolicyReconciler reconciles a ORYAccessControlPolicy object
type ORYAccessControlPolicyReconciler struct {
	KetoClient      KetoClientInterface
	KetoClientMaker KetoClientMakerFunc
	Log             logr.Logger
	otherClients    map[clientMapKey]KetoClientInterface
	client.Client
}

// +kubebuilder:rbac:groups=keto.ory.sh,resources=oryaccesscontrolpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keto.ory.sh,resources=oryaccesscontrolpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *ORYAccessControlPolicyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("oryaccesscontrolpolicy", req.NamespacedName)

	var oryPolicy ketov1alpha1.ORYAccessControlPolicy
	if err := r.Get(ctx, req.NamespacedName, &oryPolicy); err != nil {
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	ketoClient, err := r.getKetoClientForClient(oryPolicy)
	if err != nil {
		r.Log.Error(err, fmt.Sprintf(
			"keto address %s:%d is invalid",
			oryPolicy.Spec.Keto.URL,
			oryPolicy.Spec.Keto.Port,
		))
		if updateErr := r.updateReconciliationStatusError(ctx, &oryPolicy, ketov1alpha1.StatusInvalidKetoAddress, err); updateErr != nil {
			return ctrl.Result{}, updateErr
		}
		return ctrl.Result{}, nil
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if oryPolicy.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(oryPolicy.ObjectMeta.Finalizers, FinalizerName) {
			oryPolicy.ObjectMeta.Finalizers = append(oryPolicy.ObjectMeta.Finalizers, FinalizerName)
			if err := r.Update(ctx, &oryPolicy); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(oryPolicy.ObjectMeta.Finalizers, FinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := ketoClient.DeleteORYAccessControlPolicy(oryPolicy.Spec.Flavor, req.NamespacedName.Name); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			oryPolicy.ObjectMeta.Finalizers = removeString(oryPolicy.ObjectMeta.Finalizers, FinalizerName)
			if err := r.Update(ctx, &oryPolicy); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil

	}

	if oryPolicy.Generation != oryPolicy.Status.ObservedGeneration {

		if _, err := ketoClient.PutORYAccessControlPolicy(oryPolicy.Spec.Flavor, oryPolicy.ToORYAccessControlPolicyJSON()); err != nil {
			if updateErr := r.updateReconciliationStatusError(ctx, &oryPolicy, ketov1alpha1.StatusUpsertPolicyFailed, err); updateErr != nil {
				return ctrl.Result{}, updateErr
			}
			return ctrl.Result{}, nil
		}

		if err := r.ensureEmptyStatusError(ctx, &oryPolicy); err != nil {
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

func (r *ORYAccessControlPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ketov1alpha1.ORYAccessControlPolicy{}).
		Complete(r)
}

func (r *ORYAccessControlPolicyReconciler) updateReconciliationStatusError(ctx context.Context, c *ketov1alpha1.ORYAccessControlPolicy, code ketov1alpha1.StatusCode, err error) error {
	r.Log.Error(err, fmt.Sprintf("error processing client %s/%s ", c.Name, c.Namespace), "oryPolicy", "register")
	c.Status.ReconciliationError = ketov1alpha1.ReconciliationError{
		Code:        code,
		Description: err.Error(),
	}

	return r.updatePolicyStatus(ctx, c)
}

func (r *ORYAccessControlPolicyReconciler) ensureEmptyStatusError(ctx context.Context, c *ketov1alpha1.ORYAccessControlPolicy) error {
	c.Status.ReconciliationError = ketov1alpha1.ReconciliationError{}
	return r.updatePolicyStatus(ctx, c)
}

func (r *ORYAccessControlPolicyReconciler) updatePolicyStatus(ctx context.Context, c *ketov1alpha1.ORYAccessControlPolicy) error {
	c.Status.ObservedGeneration = c.Generation
	if err := r.Status().Update(ctx, c); err != nil {
		r.Log.Error(err, fmt.Sprintf("status update failed for client %s/%s ", c.Name, c.Namespace), "oryPolicy", "update status")
		return err
	}
	return nil
}

func (r *ORYAccessControlPolicyReconciler) getKetoClientForClient(oryPolicy KetoConfiger) (KetoClientInterface, error) {
	ketoConfig := oryPolicy.GetKeto()
	if ketoConfig == (ketov1alpha1.Keto{}) {
		r.Log.Info(fmt.Sprintf("using default client"))
		return r.KetoClient, nil
	}
	key := clientMapKey{
		url:      ketoConfig.URL,
		port:     ketoConfig.Port,
		endpoint: ketoConfig.Endpoint,
	}
	if c, ok := r.otherClients[key]; ok {
		return c, nil
	}
	return r.KetoClientMaker(oryPolicy)
}
