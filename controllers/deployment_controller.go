/*
Copyright 2022.

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
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientutil "kmodules.xyz/client-go/client"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("deployment", req.NamespacedName)

	logger.Info(fmt.Sprintf("**************** hello from deployment reconciler %s ****************", req.String()))
	// Ignore if Namepsace is equal to "Kube-system"

	if !IncludedNamespace(req.Namespace) {
		logger.Info(fmt.Sprintf("drop the key %s as excluded namespace", req.String()))
		return ctrl.Result{}, nil
	}

	// Getting the Deployment Object
	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deployment); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	containers, err := pushContainersToBackupRegistry(ctx, r.Client, deployment.Spec.Template.Spec.Containers)
	if err != nil {
		return ctrl.Result{}, err
	}

	initContainers, err := pushContainersToBackupRegistry(ctx, r.Client, deployment.Spec.Template.Spec.InitContainers)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, _, err = clientutil.CreateOrPatch(ctx, r.Client, deployment, func(obj client.Object, createOp bool) client.Object {
		in := obj.(*appsv1.Deployment)
		in.Spec.Template.Spec.Containers = containers
		in.Spec.Template.Spec.InitContainers = initContainers
		return in
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(obj client.Object) bool {
			return IncludedNamespace(obj.GetNamespace())
			//return IncludedNamespace(obj.GetNamespace()) && !meta_util.MustAlreadyReconciled(obj)
		})).
		Complete(r)
}
