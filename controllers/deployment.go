package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vortalbizv1 "vortal.biz/joaoneves/azdevops-operator/api/v1"
)

func labels(v *vortalbizv1.AzDevopsAgentPool, tier string) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app":             "visitors",
		"visitorssite_cr": v.Name,
		"tier":            tier,
	}
}

// ensureDeployment ensures Deployment resource presence in given namespace.
func (r *AzDevopsAgentPoolReconciler) ensureDeployment(request reconcile.Request,
	instance *vortalbizv1.AzDevopsAgentPool,
	dep *appsv1.Deployment,
	ctx context.Context,
) (*reconcile.Result, error) {

	log := log.FromContext(ctx).WithValues("AzDevopsController", dep.Name)

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	replicas := instance.Spec.Autoscaling.Min
	if *found.Spec.Replicas != replicas {
		found.Spec.Replicas = &replicas
		err = r.Update(context.TODO(), found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &ctrl.Result{}, err
		}
		// Spec updated return and requeue
		// Requeue for any reason other than an error
		return &ctrl.Result{Requeue: true}, nil
	}

	return nil, nil
}

// backendDeployment is a code for Creating Deployment
func (r *AzDevopsAgentPoolReconciler) backendDeployment(v *vortalbizv1.AzDevopsAgentPool) *appsv1.Deployment {

	labels := labels(v, "backend")
	size := int32(v.Spec.Autoscaling.Min)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hello-pod",
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v.Spec.Template,
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}
