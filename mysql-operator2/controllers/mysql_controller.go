/*
Copyright 2021.

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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cachev1alpha1 "github.com/Sher-Chowdhury/mysql-operator2/api/v1alpha1"
	pod "github.com/Sher-Chowdhury/mysql-operator2/resources/pod"
)

// MysqlReconciler reconciles a Mysql object
type MysqlReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.codingbee.net,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.codingbee.net,resources=mysqls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.codingbee.net,resources=mysqls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Mysql object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *MysqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("mysql", req.NamespacedName)

	// your logic here

	// Fetch the MySQL CR's yaml data.

	// First, create an empty variable that will hold the MySQL CR's yaml data
	instance := &cachev1alpha1.Mysql{}

	// Second, use the Client.Get method to populate the 'instance' variable
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
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
	// debugging: this should print out the mysql CR's yaml data:
	fmt.Println("mysql cr data:")
	fmt.Println(instance)

	// Define a new Pod object. This is the equivalent of writing a pod yaml file, but not applying it yet.
	var mysql_pod = pod.MysqlPodDef
	pod := mysql_pod(instance)

	// debugging: this should print out the mysql pod def that we're going to build:
	fmt.Println("New pod definition that will be used to create a new pod:")
	fmt.Println(pod)

	// Set MySQL Cr as the owner of this pod definition (so that this pod will get deleted if the CR is deleted)
	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Create an empty pod variable, which will be populated with pod data, if a match is found.
	found := &corev1.Pod{}

	// populate the "found" with a match (if a match is found)
	// if a match is not found, then we get an error, which get's captured in 'err'
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)

	// if we get an error of the time "IsNotFound" then proceed with running the "oc create" command.
	if err != nil && errors.IsNotFound(err) {

		// This create line is essentially doing the "oc create" command.
		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		// return reconcile.Result{}, nil    // I don't think this line is needed.
		fmt.Println("The following pod is now created: " + pod.Name)

	} else if err != nil {
		// if instead we get some other kind of error then raise it as an actual error.
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Mysql{}).
		Complete(r)
}
