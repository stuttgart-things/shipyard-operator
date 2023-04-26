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
	"os"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	// "github.com/hashicorp/terraform-exec/tfexec"
	shipyardv1beta1 "github.com/stuttgart-things/shipyard-operator/api/v1beta1"
)

// ShipyardTerraformReconciler reconciles a ShipyardTerraform object
type ShipyardTerraformReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const regexPatternVaultSecretPath = `.+/data/.+:.+`

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
	err := r.Get(ctx, req.NamespacedName, terraformCR)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Terraform resource not found...")
		} else {
			log.Info("Error", err)
		}
	}

	// GET VARIABLES FROM CR
	var tfVersion string = terraformCR.Spec.TerraformVersion
	var template string = terraformCR.Spec.Template
	var module []string = terraformCR.Spec.Module
	var backend []string = terraformCR.Spec.Backend
	var secrets []string = terraformCR.Spec.Secrets
	var variables []string = terraformCR.Spec.Variables
	var workingDir = "/tmp/tf/" + req.Name + "/"
	var tfInitOptions []tfexec.InitOption
	var applyOptions []tfexec.ApplyOption

	// GET MODULE PARAMETER
	moduleParameter := make(map[string]interface{})
	for _, s := range module {
		keyValue := strings.Split(s, "=")
		moduleParameter[keyValue[0]] = keyValue[1]
	}

	// CHECK FOR VAULT ENV VARS
	vaultAuthType, vaultAuthFound := VerifyVaultEnvVars()
	log.Info("⚡️ VAULT CREDENDITALS ⚡️", vaultAuthType, vaultAuthFound)

	// CONVERT ALL EXISTING SECRETS IN BACKEND+MODULE PARAMETERS
	backend = ConvertVaultSecretsInParameters(backend)
	secrets = ConvertVaultSecretsInParameters(secrets)

	// PRINT OUT CR
	fmt.Println("CR-NAME", req.Name)

	// READ + RENDER TF MODULE TEMPLATE
	moduleCallTemplate := sthingsBase.ReadFileToVariable("terraform/" + template)
	log.Info("⚡️ Rendering tf config ⚡️")
	renderedModuleCall, _ := sthingsBase.RenderTemplateInline(string(moduleCallTemplate), "missingkey=zero", "{{", "}}", moduleParameter)

	// CREATE TF FILES
	log.Info("⚡️ Creating working dir and project files ⚡️")
	sthingsBase.CreateNestedDirectoryStructure(workingDir, 0777)
	sthingsBase.StoreVariableInFile(workingDir+req.Name+".tf", string(renderedModuleCall))
	sthingsBase.StoreVariableInFile(workingDir+"terraform.tfvars", strings.Join(variables, "\n"))

	// TERRAFORM INIT
	tf := InitalizeTerraform(workingDir, tfVersion)
	log.Info("⚡️ Initalize terraform ⚡️")
	tfInitOptions = append(tfInitOptions, tfexec.Upgrade(true))

	for _, backendParameter := range backend {
		tfInitOptions = append(tfInitOptions, tfexec.BackendConfig(strings.TrimSpace(backendParameter)))
	}

	err = tf.Init(context.Background(), tfInitOptions...)

	if err != nil {
		fmt.Println("error running Init: %s", err)
	}

	log.Info("⚡️ Initalizing terraform done ⚡️")

	// TERRAFORM APPLY
	log.Info("⚡️ Applying.. ⚡️")
	for _, secret := range secrets {
		applyOptions = append(applyOptions, tfexec.Var(strings.TrimSpace(secret)))
	}

	err = tf.Apply(context.Background(), applyOptions...)

	if err != nil {
		fmt.Println("error running Apply: %s", err)
	}

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	log.Info("TF APPLY DONE!")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShipyardTerraformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipyardv1beta1.ShipyardTerraform{}).
		Complete(r)
}
