/*
Copyright 2024 yht.

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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8sappsv1 "k8s.io/api/apps/v1"
	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "example/api/v1"
)

// PodsbookReconciler reconciles a Podsbook object
type PodsbookReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=apps.ppdapi.com,resources=podsbooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.ppdapi.com,resources=podsbooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.ppdapi.com,resources=podsbooks/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Podsbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *PodsbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	podsbook := &appsv1.Podsbook{}
	deployment := &k8sappsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, podsbook)
	if err != nil {
		return ctrl.Result{}, nil
	}
	err = r.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Deployment Not Found")
			err = r.CreateDeployment(ctx, podsbook)
			if err != nil {
				r.Recorder.Event(podsbook, k8scorev1.EventTypeWarning, "FailedCreateDeployment", err.Error())
				return ctrl.Result{}, err
			}
			podsbook.Status.RealReplica = *podsbook.Spec.Replica
			err = r.Update(ctx, podsbook)
			if err != nil {
				r.Recorder.Event(podsbook, k8scorev1.EventTypeWarning, "FailedUpdateStatus", err.Error())
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, err
	}
	//binding deployment to podsbook
	if err = ctrl.SetControllerReference(podsbook, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

func (r *PodsbookReconciler) CreateDeployment(ctx context.Context, podsbook *appsv1.Podsbook) error {
	deployment := &k8sappsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: podsbook.Namespace,
			Name:      podsbook.Name,
		},
		Spec: k8sappsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(*podsbook.Spec.Replica),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": podsbook.Name,
				},
			},

			Template: k8scorev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": podsbook.Name,
					},
				},
				Spec: k8scorev1.PodSpec{
					Containers: []k8scorev1.Container{
						{
							Name:            podsbook.Name,
							Image:           *podsbook.Spec.Image,
							ImagePullPolicy: "IfNotPresent",
							Ports: []k8scorev1.ContainerPort{
								{
									Name:          podsbook.Name,
									Protocol:      k8scorev1.ProtocolSCTP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	err := r.Create(ctx, deployment)
	if err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodsbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Podsbook{}).
		Owns(&k8sappsv1.Deployment{}).
		Complete(r)
}
