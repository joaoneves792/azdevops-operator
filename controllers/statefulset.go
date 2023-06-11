package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vortalbizv1 "vortal.biz/joaoneves/azdevops-operator/api/v1"
)

func labels(v *vortalbizv1.AzDevopsAgentPool) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app":       "azdevops-agent-pool",
		"pool-name": v.Name,
	}
}

// ensureStatefulSet ensures StatefulSet resource presence in given namespace.
func (r *AzDevopsAgentPoolReconciler) ensureStatefulSet(request reconcile.Request,
	instance *vortalbizv1.AzDevopsAgentPool,
	sts *appsv1.StatefulSet,
	ctx context.Context,
) (*reconcile.Result, error) {

	log := log.FromContext(ctx).WithValues("AzDevopsController", sts.Name)

	// See if StatefulSet already exists and create if it doesn't
	found := &appsv1.StatefulSet{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      sts.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the stsloyment
		err = r.Create(context.TODO(), sts)

		if err != nil {
			// StatefulSet failed
			return &reconcile.Result{}, err
		} else {
			// StatefulSet was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the StatefulSet not existing
		return &reconcile.Result{}, err
	}

	replicas := instance.Spec.Autoscaling.Min
	if *found.Spec.Replicas != replicas {
		found.Spec.Replicas = &replicas
		err = r.Update(context.TODO(), found)
		if err != nil {
			log.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
			return &ctrl.Result{}, err
		}
		// Spec updated return and requeue
		// Requeue for any reason other than an error
		return &ctrl.Result{Requeue: true}, nil
	}

	return nil, nil
}

func (r *AzDevopsAgentPoolReconciler) poolStatefulSet(v *vortalbizv1.AzDevopsAgentPool) *appsv1.StatefulSet {

	labels := labels(v)
	size := int32(v.Spec.Autoscaling.Min)
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name,
			Namespace: v.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v.Spec.Template.Spec,
			},
		},
	}

	controllerutil.SetControllerReference(v, sts, r.Scheme)
	return sts
}
