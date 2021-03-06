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
	"log"
	"os/exec"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appv1alpha1 "github.com/jebinjeb/onprem-operator/api/v1alpha1"
)

// ChartReconciler reconciles a Chart object
type ChartReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.onprem.local,resources=charts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.onprem.local,resources=charts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.onprem.local,resources=charts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Chart object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ChartReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("chart", req.NamespacedName)

	// your logic here
	// Fetch the Chart instance
	instance := &appv1alpha1.Chart{}
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
	// Prinit Chart spec
	chartSpec := instance.Spec
	fmt.Printf("Chart spec received from cr is %s\n", chartSpec)

	err = deployHelm(chartSpec)
	if err != nil {
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChartReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Chart{}).
		Complete(r)
}

func deployHelm(chartSpec appv1alpha1.ChartSpec) error {
	out, err := exec.Command("helm", "install", "rainier-nginx", chartSpec.ChartName).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The result of helm install:  %s\n", out)
	return nil
}
