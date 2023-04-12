# shipyard-operator
// TODO(user): Add simple overview of use/purpose

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster



### Undeploy controller
UnDeploy the controller from the cluster.

### Create ShipyardTerraform CR

```
apiVersion: shipyard.sthings.tiab.ssc.sva.de/v1beta1
kind: ShipyardTerraform
metadata:
  name: shipyardterraform-sample
  labels:
    app.kubernetes.io/name: shipyardterraform
    app.kubernetes.io/part-of: shipyard-operator
    app.kubernetes.io/created-by: shipyard-operator
spec:
  variables:
    - vsphere_vm_name="shipyard-operator5"
    - vm_count=1
    - vm_num_cpus=6
    - vm_memory=8192
    - vsphere_vm_template="/LabUL/host/Cluster01/10.31.101.40/ubuntu22"
    - vsphere_vm_folder_path="phermann/rancher-things"
    - vsphere_network="/LabUL/host/Cluster01/10.31.101.41/MGMT-10.31.101"
    - vsphere_datastore="/LabUL/host/Cluster01/10.31.101.41/UL-ESX-SAS-01"
    - vsphere_resource_pool="/LabUL/host/Cluster01/Resources"
    - vsphere_datacenter="LabUL"
  module:
    - moduleName=shipyard-operator5
    - backendKey=shipyard-operator5.tfstate
    - moduleSourceUrl=https://artifacts.tiab.labda.sva.de/modules/vsphere-vm.zip
    - backendEndpoint=https://artifacts.tiab.labda.sva.de
    - backendRegion=main
    - backendBucket=vsphere-vm
    - tfProviderName=vsphere
    - tfProviderSource=hashicorp/vsphere
    - tfProviderVersion=2.3.1
    - tfVersion=1.4.4
  backend:
    - access_key=apps/data/artifacts:rootUser
    - secret_key=apps/data/artifacts:rootPassword
  secrets:
    - vsphere_user=cloud/data/vsphere:username
    - vsphere_password=cloud/data/vsphere:password
    - vsphere_server=cloud/data/vsphere:ip
    - vm_ssh_user=cloud/data/vsphere:vm_ssh_user
    - vm_ssh_password=cloud/data/vsphere:vm_ssh_password
  terraform-version: 1.4.4
  template: vsphere-vm
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.


## License

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
