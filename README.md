# stuttgart-things/shipyard-operator

shipyard-operator allows you to manage the lifecycle of terraform resources through custom resources.

## Getting Started

<details><summary>RUN SHIPYARD-OPERATOR W/ HELM</summary>



#### STABLE

#### LATEST DEV RELEASE

```
helm upgrade --install shipyard-operator \
oci://eu.gcr.io/stuttgart-things/shipyard-operator:v0.1.62 \
-n shipyard-operator-system --create-namespace
```

</details>
 
<details><summary>DEPLOY W/ ARGOCD-VAULT-PLUGIN</summary>

[argocd-vault-plugin](https://argocd-vault-plugin.readthedocs.io/en/stable/)


```
cat <<EOF > shipyard-operator.yaml
secrets:
  vault:
    name: vault
    labels:
      app.kubernetes.io/component: manager
      app.kubernetes.io/created-by: shipyard-operator
      app.kubernetes.io/instance: controller-manager
      app.kubernetes.io/part-of: shipyard-operator
    secretKVs:
      VAULT_TOKEN: <path:apps/data/vault#token>
      VAULT_NAMESPACE: <path:apps/data/vault#namespace>
      VAULT_ADDR: <path:apps/data/vault#url>
EOF

# Template/test

export VAULT_ADDR=https://<VAULT-URL>
export VAULT_NAMESPACE=<NAMESPACE>
export VAULT_TOKEN=<NAMESPACE>

helm template --values <(cat shipyard-operator.yaml) . | argocd-vault generate -

# Install/Upgrade

export VAULT_ADDR=https://<VAULT-URL>
export VAULT_NAMESPACE=<NAMESPACE>
export VAULT_TOKEN=<NAMESPACE>

helm template -n shipyard-operator-system \
--values <(cat dev2.yaml) . | argocd-vault generate - | kubectl apply -f -
```

</details>


### Undeploy controller
UnDeploy the controller from the cluster.

### Create ShipyardTerraform CR

<details><summary>EXAMPLE-VSPHERE-VM</summary>

```
apiVersion: shipyard.sthings.tiab.ssc.sva.de/v1beta1
kind: ShipyardTerraform
metadata:
  name: yacht-vm1
  labels:
    app.kubernetes.io/name: shipyardterraform
    app.kubernetes.io/part-of: shipyard-operator
    app.kubernetes.io/created-by: shipyard-operator
spec:
  variables:
    - vsphere_vm_name="yacht1"
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
    - moduleName=yacht1
    - backendKey=yacht1.tfstate
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

</details>

<details><summary>EXAMPLE-PVE-VM</summary>

```
apiVersion: shipyard.sthings.tiab.ssc.sva.de/v1beta1
kind: ShipyardTerraform
metadata:
  name: shipyardterraform-pve-sample
  labels:
    app.kubernetes.io/name: shipyardterraform
    app.kubernetes.io/part-of: shipyard-operator
    app.kubernetes.io/created-by: shipyard-operator
spec:
  variables:
    - vm_name="shipyard-operator-pve1"
    - vm_count=1
    - vm_num_cpus=6
    - vm_memory=8192
    - vm_template="u22-rke2-upi"
    - pve_network="vmbr101"
    - pve_datastore="v3700"
    - vm_disk_size="128G"
    - pve_folder_path="stuttgart-things"
    - pve_cluster_node="sthings-pve1"
  module:
    - moduleName=shipyard-operator-pve1
    - backendKey=shipyard-operator-pve1.tfstate
    - moduleSourceUrl=https://artifacts.app.sthings-pve.labul.sva.de/modules/proxmox-vm.zip
    - backendEndpoint=https://artifacts.app.sthings-pve.labul.sva.de
    - backendRegion=main
    - backendBucket=pve-vm
    - tfProviderName=proxmox
    - tfProviderSource=Telmate/proxmox
    - tfProviderVersion=2.9.14
    - tfVersion=1.4.4
  backend:
    - access_key=apps/data/artifacts:rootUser
    - secret_key=apps/data/artifacts:rootPassword
  secrets:
    - pve_api_url=cloud/data/pve:api_url
    - pve_api_user=cloud/data/pve:api_user
    - pve_api_password=cloud/data/pve:api_password
    - vm_ssh_user=cloud/data/pve:ssh_user
    - vm_ssh_password=cloud/data/pve:ssh_password
  terraform-version: 1.4.5
  template: pve-vm
```

</details>

  
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
