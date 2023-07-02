/*
Copyright 2023.

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
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"

	vortalbizv1 "vortal.biz/joaoneves/azdevops-operator/api/v1"
)

// AzDevopsAgentPoolReconciler reconciles a AzDevopsAgentPool object
type AzDevopsAgentPoolReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	PAT    string
	log    logr.Logger
}

//+kubebuilder:rbac:groups=vortal.biz,resources=azdevopsagentpools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vortal.biz,resources=azdevopsagentpools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vortal.biz,resources=azdevopsagentpools/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AzDevopsAgentPool object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AzDevopsAgentPoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := log.FromContext(ctx).WithValues("AzDevopsController", req.NamespacedName)
	r.log = log

	// Fetch the AZDevopsPool instance
	instance := &vortalbizv1.AzDevopsAgentPool{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if this StatefulSet already exists
	found := &appsv1.StatefulSet{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	desiredReplicas := int32(0)
	if err == nil {
		log.Info("Running autoscaling logic")
		desiredReplicas, err = r.autoscale(req, instance, found, ctx)
		if err != nil {
			log.Error(err, "Failed to calculate desired number of replicas")
			desiredReplicas = 0
		}
		log.Info("Autoscaler has finished", "desiredReplicas", strconv.Itoa(int(desiredReplicas)))

	}

	var result *reconcile.Result
	sts := r.poolStatefulSet(instance)
	result, err = r.ensureStatefulSet(req, instance, sts, ctx, desiredReplicas)
	if result != nil {
		log.Error(err, "StatefulSet pool Not ready")
		return *result, err
	}

	// StatefulSet and Service already exists - don't requeue
	log.Info("Reconcile OK: StatefulSet and service already exists",
		"StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AzDevopsAgentPoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vortalbizv1.AzDevopsAgentPool{}).
		Complete(r)
}
