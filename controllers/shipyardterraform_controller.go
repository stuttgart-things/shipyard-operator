/*
Copyright 2023 patrick hermann.

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

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	shipyardv1beta1 "github.com/stuttgart-things/shipyard-operator/api/v1beta1"
)

// ShipyardTerraformReconciler reconciles a ShipyardTerraform object
type ShipyardTerraformReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=shipyard.sthings.tiab.ssc.sva.de,resources=shipyardterraforms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=shipyard.sthings.tiab.ssc.sva.de,resources=shipyardterraforms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=shipyard.sthings.tiab.ssc.sva.de,resources=shipyardterraforms/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ShipyardTerraform object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ShipyardTerraformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := ctrllog.FromContext(ctx)
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	terraformCR := &shipyardv1beta1.ShipyardTerraform{}
	var tfVersion string = terraformCR.Spec.TerraformVersion
	var template string = terraformCR.Spec.Template

	err := r.Get(ctx, req.NamespacedName, terraformCR)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Terraform resource not found...")
		} else {
			log.Info("Error", err)
		}
	}

	fmt.Println("CR-NAME", req.Name)
	fmt.Println("TEMPLATE", template)
	fmt.Println("TF-VERSION", tfVersion)

	b := sthingsBase.ReadFileToVariable("terraform/vsphere-vm")
	fmt.Println(string(b))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShipyardTerraformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipyardv1beta1.ShipyardTerraform{}).
		Complete(r)
}
