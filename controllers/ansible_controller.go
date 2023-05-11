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
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	shipyardv1beta1 "github.com/stuttgart-things/shipyard-operator/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ = kubernetes.NewForConfig(config)
)

type AnsibleJob struct {
	Name      string
	Namespace string
	Image     string
}

const AnsibleJobTemplate = `
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    jobgroup: ansible
    app: shipyard-operator
spec:
  template:
    metadata:
      name: {{ .Name }}
      namespace: {{ .Namespace }}
      labels:
        jobgroup: ansible
        app: shipyard-operator
    spec:
      volumes:
        - name: ansible-templates
          configMap:
            name: ansible-templates
        - name: workdir
          emptyDir:
            medium: Memory
      securityContext:
        runAsUser: 65532
        fsGroup: 65532
      initContainers:
        - name: prepare-workdir
          image: {{ .Image }}
          args: ["cp /ansible/inventory /workdir/inventory && cat /workdir/inventory"]
          command:
            - sh
            - -c
          volumeMounts:
            - name: ansible-templates
              mountPath: /ansible/
            - name: workdir
              mountPath: /workdir/
      containers:
        - name: execute-ansible
          envFrom:
            - secretRef:
                name: vault
          image: {{ .Image }}
          args: ["ansible-playbook -i /workdir/inventory /ansible/play.yaml -vv -e inventory_path=/workdir/inventory"]
          command:
            - sh
            - -c
          env:
            - name: ANSIBLE_ROLES_PATH
              value: "/home/nonroot/.ansible/roles"
            - name: ANSIBLE_UNSAFE_WRITES
              value: "1"
            - name: ANSIBLE_HOST_KEY_CHECKING
              value: "False"
          volumeMounts:
            - name: ansible-templates
              mountPath: /ansible/
            - name: workdir
              mountPath: /workdir/
      restartPolicy: Never
`

// AnsibleReconciler reconciles a Ansible object
type AnsibleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

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

	log := ctrllog.FromContext(ctx)
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	// createJob
	job := AnsibleJob{
		Name:      "baseos",
		Namespace: "shipyard-operator-system",
		Image:     "eu.gcr.io/stuttgart-things/sthings-ansible:7.5.0-3",
	}

	fmt.Println(job)

	tmpl, err := template.New("pipelinerun").Parse(AnsibleJobTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, job)

	if err != nil {
		fmt.Println("execution: %s", err)
	}

	fmt.Println("RENDERED JOB")
	fmt.Println(buf.String())

	if watchJobExecution("countdown") {
		fmt.Println("LETS HOPE FOR NO RESTART")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AnsibleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipyardv1beta1.Ansible{}).
		Complete(r)
}

func watchJobExecution(jobName string) (jobSuccessful bool) {
	timeOut := int64(60)
	watcher, _ := clientset.BatchV1().Jobs("").Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeOut})

Jobs:
	for event := range watcher.ResultChan() {
		item := event.Object.(*batchv1.Job)

		switch event.Type {
		case watch.Modified:

			fmt.Println(item.GetName)
			if item.Name == jobName {

				fmt.Println(item.Status.Active)

				if item.Status.Active == 0 {
					fmt.Println("JOB IS FINISHED!")
					return true
					break Jobs
				}

			}

		case watch.Bookmark:
		case watch.Error:
		case watch.Deleted:
		case watch.Added:

		}
	}

	fmt.Println("END OF WATCHING ANSIBLE EXECUTION")
	return jobSuccessful
}
