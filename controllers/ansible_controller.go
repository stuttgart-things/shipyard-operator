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

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	shipyardv1beta1 "github.com/stuttgart-things/shipyard-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// AnsibleReconciler reconciles a Ansible object
type AnsibleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var workingDir = "/tmp/.ansible/tmp"

//+kubebuilder:rbac:groups=shipyard.sthings.tiab.ssc.sva.de,resources=ansibles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=shipyard.sthings.tiab.ssc.sva.de,resources=ansibles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=shipyard.sthings.tiab.ssc.sva.de,resources=ansibles/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ansible object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AnsibleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TEST BLOCK BEGIN

	fmt.Println("Hello Ansible2!")
	log := ctrllog.FromContext(ctx)
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	log.Info("⚡️ Creating working dir and project files ⚡️")
	sthingsBase.CreateNestedDirectoryStructure(workingDir, 0777)

	os.Setenv("ANSIBLE_LOCAL_TEMP", workingDir)
	os.Setenv("ANSIBLE_REMOTE_TMP", workingDir)
	os.Setenv("ANSIBLE_HOST_KEY_CHECKING", "False")

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User: "sthings",
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: "elastic-vm1.tiab.labda.sva.de,",
		ExtraVars: map[string]interface{}{
			"ansible_password":   "Atlan7is",
			"ansible_remote_tmp": workingDir,
		},
	}

	ansiblePlaybookPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become:        false,
		AskBecomePass: false,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  []string{"terraform/play"},
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		Options:                    ansiblePlaybookOptions,
		Binary:                     "/venv/bin/ansible-playbook",
		Exec: execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithTransformers(
				results.Prepend("Go-ansible example with become"),
			),
		),
	}

	err := playbook.Run(context.TODO())
	if err != nil {
		panic(err)
	}

	// TEST BLOCK END
	// ANSIBLE_ROLES_PATH

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AnsibleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipyardv1beta1.Ansible{}).
		Complete(r)
}
